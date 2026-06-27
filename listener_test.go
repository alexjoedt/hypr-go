package hypr

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestListen(t *testing.T) {
	t.Run("dials the event socket under the instance signature", func(t *testing.T) {
		runtimeDir := t.TempDir()
		sig := "sig"

		dir := filepath.Join(runtimeDir, "hypr", sig)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		path := filepath.Join(dir, eventSocketName)
		if len(path) > 100 {
			t.Skipf("socket path %q is too long for sun_path", path)
		}

		ln, err := net.Listen("unix", path)
		if err != nil {
			t.Fatalf("listen: %v", err)
		}
		defer ln.Close()
		go func() {
			for {
				conn, err := ln.Accept()
				if err != nil {
					return
				}
				go holdOpen(conn)
			}
		}()

		t.Setenv(xdgRuntimeDir, runtimeDir)
		t.Setenv(hyprlandInstanceSignature, sig)

		l, err := Listen(context.Background(), WithLogger(discardLogger()))
		if err != nil {
			t.Fatalf("Listen: %v", err)
		}
		l.Close()
	})

	t.Run("missing signature", func(t *testing.T) {
		t.Setenv(hyprlandInstanceSignature, "")

		if _, err := Listen(context.Background()); !errors.Is(err, ErrNoHyprlandSocket) {
			t.Errorf("Listen error = %v, want ErrNoHyprlandSocket", err)
		}
	})

	t.Run("nothing listening", func(t *testing.T) {
		_, err := Listen(context.Background(),
			WithSocketDir(t.TempDir()),
			WithDialTimeout(time.Second),
			WithLogger(discardLogger()),
		)
		if err == nil {
			t.Fatal("Listen succeeded against an empty socket dir")
		}
		if errors.Is(err, ErrNoHyprlandSocket) {
			t.Errorf("Listen error = %v, want a dial error", err)
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := Listen(ctx,
			WithSocketDir(t.TempDir()),
			WithDialer(&fakeDialer{}),
			WithLogger(discardLogger()),
		)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Listen error = %v, want context.Canceled", err)
		}
	})

	t.Run("dials the configured directory", func(t *testing.T) {
		client, server := net.Pipe()
		defer server.Close()

		d := &fakeDialer{dial: func(context.Context) (net.Conn, error) { return client, nil }}

		l, err := Listen(context.Background(),
			WithDialer(d),
			WithSocketDir("/sockets"),
			WithLogger(discardLogger()),
		)
		if err != nil {
			t.Fatalf("Listen: %v", err)
		}
		defer l.Close()

		calls := d.dialed()
		if len(calls) != 1 {
			t.Fatalf("dialed %d times, want 1", len(calls))
		}
		if calls[0].network != "unix" {
			t.Errorf("network = %q, want unix", calls[0].network)
		}
		if want := "/sockets/" + eventSocketName; calls[0].address != want {
			t.Errorf("address = %q, want %q", calls[0].address, want)
		}
	})
}

func TestOn(t *testing.T) {
	t.Parallel()

	t.Run("typed handler sees only its type", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t)

		var got []*WorkspaceEvent
		On(l, func(e *WorkspaceEvent) { got = append(got, e) })

		l.deliver([]byte("workspace>>2\n"), yieldAll)
		l.deliver([]byte("openlayer>>waybar\n"), yieldAll)
		l.deliver([]byte("notanevent>>x\n"), yieldAll)

		if len(got) != 1 || got[0].WorkspaceName != "2" {
			t.Errorf("typed handler got %+v, want one WorkspaceEvent 2", got)
		}
	})

	t.Run("Event handler receives everything", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t)

		var got []string
		On(l, func(e Event) { got = append(got, e.Name()) })

		l.deliver([]byte("workspace>>2\n"), yieldAll)
		l.deliver([]byte("openlayer>>waybar\n"), yieldAll)
		l.deliver([]byte("notanevent>>x\n"), yieldAll)

		want := []string{"workspace", "openlayer", "notanevent"}
		if !equalStrings(got, want) {
			t.Errorf("wildcard handler saw %v, want %v", got, want)
		}
	})

	t.Run("UnknownEvent handler", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t)

		var got []*UnknownEvent
		On(l, func(e *UnknownEvent) { got = append(got, e) })

		l.deliver([]byte("workspace>>2\n"), yieldAll)
		l.deliver([]byte("notanevent>>x\n"), yieldAll)

		if len(got) != 1 || got[0].Name() != "notanevent" {
			t.Errorf("got %+v, want one UnknownEvent notanevent", got)
		}
	})

	t.Run("handlers run in registration order", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t)

		var calls []string
		for _, id := range []string{"a", "b", "c"} {
			On(l, func(*WorkspaceEvent) { calls = append(calls, id) })
		}

		l.deliver([]byte("workspace>>2\n"), yieldAll)

		if !equalStrings(calls, []string{"a", "b", "c"}) {
			t.Errorf("handlers fired as %v, want a b c", calls)
		}
	})

	t.Run("unsubscribe removes exactly one handler", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t)

		var kept, dropped int
		On(l, func(*WorkspaceEvent) { kept++ })
		unsubscribe := On(l, func(*WorkspaceEvent) { dropped++ })

		unsubscribe()
		unsubscribe() // must stay safe when called twice

		l.deliver([]byte("workspace>>2\n"), yieldAll)

		if kept != 1 {
			t.Errorf("remaining handler fired %d times, want 1", kept)
		}
		if dropped != 0 {
			t.Errorf("unsubscribed handler fired %d times, want 0", dropped)
		}
	})

	t.Run("blank and malformed lines reach no handler", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t)

		var got []Event
		On(l, func(e Event) { got = append(got, e) })

		for _, line := range []string{"\n", "   \n", "", "garbage\n"} {
			l.deliver([]byte(line), yieldAll)
		}

		if len(got) != 0 {
			t.Fatalf("handler fired %d times on unparseable input, want 0", len(got))
		}
	})
}

// TestHandlerPanic covers a handler taking down the read loop. Handlers are
// somebody else's code and run inline, so one bad callback must not cost every
// other handler its events.
func TestHandlerPanic(t *testing.T) {
	t.Parallel()

	t.Run("other handlers still run", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t)

		var before, after int
		On(l, func(*WorkspaceEvent) { before++ })
		On(l, func(*WorkspaceEvent) { panic("handler exploded") })
		On(l, func(*WorkspaceEvent) { after++ })

		l.deliver([]byte("workspace>>2\n"), yieldAll)

		if before != 1 || after != 1 {
			t.Errorf("surviving handlers fired %d and %d times, want 1 and 1", before, after)
		}
	})

	t.Run("a nil panic value is recovered too", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t)
		On(l, func(Event) { panic(nil) })

		l.deliver([]byte("workspace>>2\n"), yieldAll)
	})

	t.Run("the loop keeps reading", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t)

		events := make(chan Event, 4)
		On(l, func(e Event) {
			if e.Name() == EventWorkspace {
				panic("boom")
			}
			events <- e
		})

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errc := runAsync(ctx, l)
		defer func() { cancel(); <-errc }()

		writeLines(t, server, "workspace>>2", "openlayer>>waybar")

		got := waitEvent(t, events)
		if got.Name() != EventOpenLayer {
			t.Errorf("got %q, want the event after the panicking one", got.Name())
		}

		select {
		case err := <-errc:
			t.Fatalf("Run returned %v after a handler panicked", err)
		default:
		}
	})
}

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("delivers events in order", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t)
		events := collect(l, 8)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errc := runAsync(ctx, l)

		writeLines(t, server,
			"workspace>>2",
			"activewindow>>firefox,Home",
			"openwindow>>addr,2,firefox,a,b",
		)

		want := []Event{
			&WorkspaceEvent{WorkspaceName: "2"},
			&ActiveWindowEvent{Class: "firefox", Title: "Home"},
			&OpenWindowEvent{Address: "addr", WorkspaceName: "2", Class: "firefox", Title: "a,b"},
		}
		for i, w := range want {
			got := waitEvent(t, events)
			if got.Name() != w.Name() {
				t.Errorf("event %d: name = %q, want %q", i, got.Name(), w.Name())
			}
		}

		cancel()
		if err := waitErr(t, errc); !errors.Is(err, context.Canceled) {
			t.Errorf("Run returned %v, want context.Canceled", err)
		}
		if err := l.Err(); !errors.Is(err, context.Canceled) {
			t.Errorf("Err() = %v, want context.Canceled", err)
		}
	})

	t.Run("returns nil when the peer closes", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errc := runAsync(ctx, l)

		writeLines(t, server, "workspace>>2")
		server.Close()

		if err := waitErr(t, errc); err != nil {
			t.Errorf("Run returned %v, want nil", err)
		}
		if err := l.Err(); err != nil {
			t.Errorf("Err() = %v, want nil", err)
		}
	})

	t.Run("delivers a final line with no trailing newline", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t)
		events := collect(l, 1)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errc := runAsync(ctx, l)

		if _, err := server.Write([]byte("workspace>>7")); err != nil {
			t.Fatalf("write: %v", err)
		}
		server.Close()

		got := waitEvent(t, events)
		if ws, ok := got.(*WorkspaceEvent); !ok || ws.WorkspaceName != "7" {
			t.Errorf("got %+v, want WorkspaceEvent 7", got)
		}
		_ = waitErr(t, errc)
	})

	t.Run("keeps reading past a read timeout", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t, WithReadTimeout(50*time.Millisecond))
		events := collect(l, 1)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errc := runAsync(ctx, l)

		// Stay quiet for several read deadlines, the way an idle desktop does.
		time.Sleep(250 * time.Millisecond)

		select {
		case err := <-errc:
			t.Fatalf("Run returned %v while idle, want it still reading", err)
		default:
		}

		writeLines(t, server, "workspace>>9")

		got := waitEvent(t, events)
		if ws, ok := got.(*WorkspaceEvent); !ok || ws.WorkspaceName != "9" {
			t.Errorf("got %+v, want WorkspaceEvent 9", got)
		}
	})

	t.Run("resumes a line split by a read timeout", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t, WithReadTimeout(50*time.Millisecond))
		events := collect(l, 1)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errc := runAsync(ctx, l)
		defer func() { cancel(); <-errc }()

		if _, err := server.Write([]byte("workspacev2>>4,")); err != nil {
			t.Fatalf("write: %v", err)
		}
		time.Sleep(150 * time.Millisecond) // let the read deadline fire mid-event
		if _, err := server.Write([]byte("four\n")); err != nil {
			t.Fatalf("write: %v", err)
		}

		got := waitEvent(t, events)
		ws, ok := got.(*WorkspaceV2Event)
		if !ok {
			t.Fatalf("got %T, want *WorkspaceV2Event", got)
		}
		if ws.WorkspaceID != 4 || ws.WorkspaceName != "four" {
			t.Errorf("got %+v, want ID 4 and name four", ws)
		}
	})

	t.Run("cancelling returns promptly", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t, WithReadTimeout(time.Hour))

		ctx, cancel := context.WithCancel(context.Background())
		errc := runAsync(ctx, l)

		time.Sleep(50 * time.Millisecond) // let the read block
		cancel()

		start := time.Now()
		err := waitErr(t, errc)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Run returned %v, want context.Canceled", err)
		}
		if elapsed := time.Since(start); elapsed > time.Second {
			t.Errorf("Run took %v to notice cancellation", elapsed)
		}
	})

	t.Run("Close unblocks a running loop", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t, WithReadTimeout(time.Hour))

		errc := runAsync(context.Background(), l)

		time.Sleep(50 * time.Millisecond) // let the read block
		if err := l.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}

		if err := waitErr(t, errc); !errors.Is(err, ErrClosed) {
			t.Errorf("Run returned %v, want ErrClosed", err)
		}
	})

	t.Run("after Close", func(t *testing.T) {
		t.Parallel()

		l, _ := newTestListener(t)
		if err := l.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}

		if err := l.Run(context.Background()); !errors.Is(err, ErrClosed) {
			t.Errorf("Run returned %v, want ErrClosed", err)
		}
		for range l.Events(context.Background()) {
			t.Fatal("Events yielded after Close")
		}
		if err := l.Err(); !errors.Is(err, ErrClosed) {
			t.Errorf("Err() = %v, want ErrClosed", err)
		}
	})

	t.Run("second loop while one is active", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t)
		events := collect(l, 1)

		ctx, cancel := context.WithCancel(context.Background())
		errc := runAsync(ctx, l)
		defer func() { cancel(); <-errc }()

		// Make sure the first loop has claimed the listener.
		writeLines(t, server, "workspace>>2")
		waitEvent(t, events)

		if err := l.Run(ctx); !errors.Is(err, ErrListening) {
			t.Errorf("second Run returned %v, want ErrListening", err)
		}
		for range l.Events(ctx) {
			t.Fatal("second Events yielded")
		}
		if err := l.Err(); !errors.Is(err, ErrListening) {
			t.Errorf("Err() = %v, want ErrListening", err)
		}

		// The first loop is still alive.
		writeLines(t, server, "workspace>>3")
		if got := waitEvent(t, events); got.Name() != EventWorkspace {
			t.Errorf("first loop got %q after the second was refused", got.Name())
		}
	})

	t.Run("blank lines reach no handler", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t)
		events := collect(l, 4)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errc := runAsync(ctx, l)
		defer func() { cancel(); <-errc }()

		// The captured stream is full of these, and each one used to hand every
		// wildcard handler a nil Event.
		writeLines(t, server, "", "", "workspace>>2", "")

		got := waitEvent(t, events)
		if got == nil {
			t.Fatal("handler received a nil Event")
		}
		if got.Name() != EventWorkspace {
			t.Errorf("first event = %q, want workspace", got.Name())
		}
	})
}

func TestEvents(t *testing.T) {
	t.Parallel()

	t.Run("yields after the handlers ran", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t)

		var order []string
		On(l, func(e *WorkspaceEvent) { order = append(order, "handler "+e.WorkspaceName) })

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			for _, line := range []string{"workspace>>1\n", "workspace>>2\n"} {
				if _, err := server.Write([]byte(line)); err != nil {
					return
				}
			}
			server.Close()
		}()

		for e := range l.Events(ctx) {
			order = append(order, "range "+e.(*WorkspaceEvent).WorkspaceName)
		}

		want := []string{"handler 1", "range 1", "handler 2", "range 2"}
		if !equalStrings(order, want) {
			t.Errorf("order = %v, want %v", order, want)
		}
		if err := l.Err(); err != nil {
			t.Errorf("Err() = %v, want nil on EOF", err)
		}
	})

	t.Run("breaking out stops cleanly", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t, WithReadTimeout(time.Hour))

		writeLinesAsync(server, "workspace>>1", "workspace>>2", "workspace>>3")

		var seen int
		for range l.Events(context.Background()) {
			seen++
			break
		}

		if seen != 1 {
			t.Errorf("saw %d events, want 1", seen)
		}
		if err := l.Err(); err != nil {
			t.Errorf("Err() = %v, want nil after break", err)
		}

		// The loop is free again: the remaining lines are still readable.
		for e := range l.Events(context.Background()) {
			if e.(*WorkspaceEvent).WorkspaceName == "3" {
				break
			}
		}
	})

	t.Run("yields UnknownEvent", func(t *testing.T) {
		t.Parallel()

		l, server := newTestListener(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		writeLinesAsync(server, "notanevent>>2")

		for e := range l.Events(ctx) {
			u, ok := e.(*UnknownEvent)
			if !ok {
				t.Fatalf("got %T, want *UnknownEvent", e)
			}
			if u.Name() != "notanevent" || u.Data() != "2" {
				t.Errorf("UnknownEvent = %+v", u)
			}
			break
		}
	})
}

// TestCloseIsIdempotent also covers a zero Listener, whose socket is nil.
func TestCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	l, _ := newTestListener(t)
	if err := l.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	if err := l.Close(); err != nil {
		t.Errorf("second Close: %v", err)
	}

	var zero Listener
	if err := zero.Close(); err != nil {
		t.Errorf("zero Listener Close: %v", err)
	}
}

// TestOnConcurrent is the race detector's job: it subscribes, unsubscribes and
// delivers at the same time.
func TestOnConcurrent(t *testing.T) {
	t.Parallel()

	l, server := newTestListener(t)

	ctx, cancel := context.WithCancel(context.Background())
	errc := runAsync(ctx, l)

	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 32 {
				unsubscribe := On(l, func(*WorkspaceEvent) {})
				unsubscribe()
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 64 {
			if _, err := server.Write([]byte("workspace>>2\n")); err != nil {
				return
			}
		}
	}()

	wg.Wait()
	cancel()
	<-errc
}

// TestLoopLeavesNoGoroutine covers the watcher goroutine that used to block on
// ctx.Done() forever whenever the loop returned first.
func TestLoopLeavesNoGoroutine(t *testing.T) {
	l, server := newTestListener(t)

	before := runtime.NumGoroutine()

	// Background context, so a leaked watcher never gets released.
	errc := runAsync(context.Background(), l)
	server.Close()
	_ = waitErr(t, errc)

	for range 100 {
		if runtime.NumGoroutine() <= before {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Errorf("goroutines went from %d to %d after Run returned", before, runtime.NumGoroutine())
}

// TestBreakLeavesNoGoroutine is the iterator's version: breaking out of the
// range must release the watcher too.
func TestBreakLeavesNoGoroutine(t *testing.T) {
	l, server := newTestListener(t)

	before := runtime.NumGoroutine()

	writeLinesAsync(server, "workspace>>1")
	for range l.Events(context.Background()) {
		break
	}

	for range 100 {
		if runtime.NumGoroutine() <= before {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Errorf("goroutines went from %d to %d after breaking out of Events", before, runtime.NumGoroutine())
}

func yieldAll(Event) bool { return true }
