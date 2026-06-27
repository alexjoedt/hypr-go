package hypr

import (
	"bufio"
	"context"
	"errors"
	"io"
	"iter"
	"log/slog"
	"net"
	"os"
	"runtime/debug"
	"slices"
	"sync"
	"time"
)

var (
	// ErrListening is reported when Run or Events starts while another loop
	// on the same Listener is still active.
	ErrListening = errors.New("listener already running")
	// ErrClosed is reported by a loop that ended because Close was called, and
	// by any loop started after it.
	ErrClosed = errors.New("listener closed")
)

// Listener owns the event socket. It is the only stateful type in the package:
// Listen dials, Run or Events read, Close hangs up.
type Listener struct {
	conn        net.Conn
	readTimeout time.Duration
	log         *slog.Logger

	mu      sync.Mutex
	running bool
	closed  bool
	err     error

	handlersMu sync.Mutex
	handlers   []subscription
	nextID     int
}

type subscription struct {
	id int
	fn func(Event)
}

// Listen dials the event socket. ctx bounds the dial only; the loop takes its
// own context in Run or Events.
func Listen(ctx context.Context, opts ...Option) (*Listener, error) {
	cfg, err := resolve(opts)
	if err != nil {
		return nil, err
	}

	dialCtx, cancel := context.WithTimeout(ctx, cfg.dialTimeout)
	defer cancel()

	conn, err := cfg.dialer.DialContext(dialCtx, "unix", cfg.eventSocket())
	if err != nil {
		return nil, err
	}

	return &Listener{
		conn:        conn,
		readTimeout: cfg.readTimeout,
		log:         cfg.log,
	}, nil
}

// On registers handler for events of type E, or for every event when E is
// Event itself. Handlers run inline on the read loop, in registration order; a
// panic is recovered and logged and the loop continues. The returned function
// removes the handler and is safe to call more than once.
func On[E Event](l *Listener, handler func(E)) (unsubscribe func()) {
	l.handlersMu.Lock()
	id := l.nextID
	l.nextID++
	l.handlers = append(l.handlers, subscription{id: id, fn: func(e Event) {
		if x, ok := e.(E); ok {
			handler(x)
		}
	}})
	l.handlersMu.Unlock()

	return func() {
		l.handlersMu.Lock()
		defer l.handlersMu.Unlock()
		l.handlers = slices.DeleteFunc(l.handlers, func(h subscription) bool { return h.id == id })
	}
}

// Run drives the read loop and fires handlers until ctx is cancelled, Close is
// called or the socket ends. It is Events without a body, followed by Err.
func (l *Listener) Run(ctx context.Context) error {
	for range l.Events(ctx) {
		continue
	}

	return l.Err()
}

// Events drives the read loop and yields every event after the handlers ran
// for it. The sequence ends when ctx is cancelled, Close is called, the socket
// ends or the caller breaks out; Err says why.
func (l *Listener) Events(ctx context.Context) iter.Seq[Event] {
	return func(yield func(Event) bool) {
		if err := l.begin(); err != nil {
			return
		}

		err := l.loop(ctx, yield)
		l.end(err)
	}
}

// Err reports why the last Run or Events stopped: nil on a clean EOF or when
// the caller broke out of the range, ctx.Err() on cancel, ErrClosed after
// Close, ErrListening for a loop that never started, else the read error.
func (l *Listener) Err() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.err
}

// Close closes the socket and unblocks a running loop. It is idempotent.
func (l *Listener) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return nil
	}
	l.closed = true

	if l.conn == nil {
		return nil
	}

	return l.conn.Close()
}

// begin claims the loop. The failure is recorded so Err reports it after an
// Events range that yielded nothing.
func (l *Listener) begin() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch {
	case l.closed:
		l.err = ErrClosed
	case l.running:
		l.err = ErrListening
	default:
		l.running = true
		l.err = nil
	}

	return l.err
}

// end releases the loop and records why it stopped. A read error after Close
// is reported as ErrClosed, whatever the socket said.
func (l *Listener) end(err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed && err != nil {
		err = ErrClosed
	}
	l.running = false
	l.err = err
}

// loop reads lines until ctx is cancelled, the socket ends or yield returns
// false. Ported from the previous Conn.Listen: the partial-line carry and the
// EOF flush keep a slow or abrupt peer from losing an event.
func (l *Listener) loop(ctx context.Context, yield func(Event) bool) error {
	// Unblock a pending read when the caller cancels. done keeps this goroutine
	// from outliving the loop when it returns on its own.
	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			_ = l.conn.SetDeadline(time.Now())
		case <-done:
		}
	}()

	reader := bufio.NewReader(l.conn)

	// carry holds a partial event whose read timed out before the newline
	// arrived, so the next read can finish it instead of losing sync.
	var carry []byte

	for ctx.Err() == nil {
		_ = l.conn.SetReadDeadline(time.Now().Add(l.readTimeout))

		line, err := reader.ReadBytes('\n')
		if len(carry) > 0 {
			line = append(carry, line...)
			carry = nil
		}

		if err == nil {
			if !l.deliver(line, yield) {
				return nil
			}
			continue
		}

		switch {
		case ctx.Err() != nil:
			return ctx.Err()
		case errors.Is(err, os.ErrDeadlineExceeded):
			carry = line
			continue
		case errors.Is(err, io.EOF):
			// A final event can arrive without its trailing newline.
			if len(line) > 0 {
				l.deliver(line, yield)
			}
			return nil
		default:
			return err
		}
	}

	return ctx.Err()
}

// deliver parses one raw line, runs the handlers and yields the event. It
// reports false when the caller wants no more events.
func (l *Listener) deliver(line []byte, yield func(Event) bool) bool {
	event, err := Parse(line)
	if err != nil {
		// Blank lines between bursts are normal; they carry no event.
		return true
	}

	l.log.Debug("event received", "event", event.Name(), "raw", string(line))

	l.handlersMu.Lock()
	handlers := slices.Clone(l.handlers)
	l.handlersMu.Unlock()

	for _, h := range handlers {
		l.invoke(h.fn, event)
	}

	return yield(event)
}

// invoke calls one handler. Handlers run inline on the read loop, so a panic in
// somebody else's callback would otherwise take down the whole listener and
// every other handler with it.
func (l *Listener) invoke(fn func(Event), event Event) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		l.log.Error("event handler panicked",
			"event", event.Name(),
			"panic", r,
			"stack", string(debug.Stack()),
		)
	}()

	fn(event)
}
