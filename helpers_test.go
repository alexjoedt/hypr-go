package hypr

import (
	"context"
	"io"
	"log/slog"
	"net"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// discardLogger keeps debug output out of test logs.
func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type dialCall struct {
	network string
	address string
}

// fakeDialer hands out connections a test prepared and records what it was
// asked to dial.
type fakeDialer struct {
	dial func(ctx context.Context) (net.Conn, error)

	mu    sync.Mutex
	calls []dialCall
}

func (d *fakeDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	d.mu.Lock()
	d.calls = append(d.calls, dialCall{network: network, address: address})
	d.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return d.dial(ctx)
}

func (d *fakeDialer) dialed() []dialCall {
	d.mu.Lock()
	defer d.mu.Unlock()

	return append([]dialCall(nil), d.calls...)
}

// newTestListener returns a Listener whose event socket is one end of a
// net.Pipe, plus the other end for the test to write events into.
func newTestListener(t *testing.T, opts ...Option) (*Listener, net.Conn) {
	t.Helper()

	client, server := net.Pipe()
	t.Cleanup(func() { server.Close() })

	d := &fakeDialer{
		dial: func(context.Context) (net.Conn, error) { return client, nil },
	}

	base := []Option{
		WithDialer(d),
		WithSocketDir("/test"),
		WithReadTimeout(2 * time.Second),
		WithLogger(discardLogger()),
	}

	l, err := Listen(context.Background(), append(base, opts...)...)
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	t.Cleanup(func() { l.Close() })

	return l, server
}

// runAsync runs Run in the background and returns a channel carrying its
// return value.
func runAsync(ctx context.Context, l *Listener) <-chan error {
	errc := make(chan error, 1)
	go func() { errc <- l.Run(ctx) }()

	return errc
}

// collect registers a wildcard handler that forwards every event into a
// buffered channel.
func collect(l *Listener, size int) <-chan Event {
	events := make(chan Event, size)
	On(l, func(e Event) { events <- e })

	return events
}

// writeLines sends raw event lines, each terminated with a newline, the way
// Hyprland does.
func writeLines(t *testing.T, w net.Conn, lines ...string) {
	t.Helper()

	for _, line := range lines {
		if _, err := w.Write([]byte(line + "\n")); err != nil {
			t.Fatalf("write %q: %v", line, err)
		}
	}
}

// writeLinesAsync is writeLines from a goroutine, for a test whose reader is
// the test goroutine itself. A net.Pipe write blocks until it is read, so the
// two cannot share one goroutine. Errors are ignored: the reader side asserts.
func writeLinesAsync(w net.Conn, lines ...string) {
	go func() {
		for _, line := range lines {
			if _, err := w.Write([]byte(line + "\n")); err != nil {
				return
			}
		}
	}()
}

// waitEvent waits for one event, failing the test if none arrives.
func waitEvent(t *testing.T, ch <-chan Event) Event {
	t.Helper()

	select {
	case e := <-ch:
		return e
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for an event")
		return nil
	}
}

// waitErr waits for Run to return.
func waitErr(t *testing.T, errc <-chan error) error {
	t.Helper()

	select {
	case err := <-errc:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for Run to return")
		return nil
	}
}

// socketDirWith starts a unix listener for each named socket inside a fresh
// temp dir and returns the dir, ready to hand to WithSocketDir.
func socketDirWith(t *testing.T, handlers map[string]func(net.Conn)) string {
	t.Helper()

	dir := t.TempDir()

	for name, handle := range handlers {
		path := filepath.Join(dir, name)
		// sun_path is 108 bytes on Linux; a long TMPDIR would blow the limit.
		if len(path) > 100 {
			t.Skipf("socket path %q is too long for sun_path", path)
		}

		ln, err := net.Listen("unix", path)
		if err != nil {
			t.Fatalf("listen %s: %v", path, err)
		}
		t.Cleanup(func() { ln.Close() })

		go func(handle func(net.Conn)) {
			for {
				conn, err := ln.Accept()
				if err != nil {
					return
				}
				go func() {
					defer conn.Close()
					handle(conn)
				}()
			}
		}(handle)
	}

	return dir
}

// holdOpen keeps a connection open until the test tears the listener down. It
// stands in for the event socket when a test only cares about commands.
func holdOpen(conn net.Conn) {
	_, _ = io.Copy(io.Discard, conn)
}
