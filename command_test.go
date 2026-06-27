package hypr

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRawCommand(t *testing.T) {
	t.Parallel()

	if got := RawCommand("hl.dsp.no_op").Command(); got != "hl.dsp.no_op" {
		t.Errorf("Command() = %q, want hl.dsp.no_op", got)
	}
}

func TestParseReply(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		reply    string
		ok       bool
		errors   []string
		warnings []string
		msg      string
	}{
		{name: "ok", reply: "ok", ok: true},
		{name: "uppercase", reply: "OK", ok: true},
		{name: "padded", reply: " ok\n", ok: true},
		{
			name:   "one error line",
			reply:  "error: attempt to call a nil value (field 'nonexistent')",
			errors: []string{"attempt to call a nil value (field 'nonexistent')"},
			msg:    "req: attempt to call a nil value (field 'nonexistent')",
		},
		{
			name:     "one warning line",
			reply:    "warning: window not found",
			warnings: []string{"window not found"},
			msg:      "req: window not found",
		},
		{
			name:     "mixed",
			reply:    "warning: first\nerror: second\nwarning: third",
			errors:   []string{"second"},
			warnings: []string{"first", "third"},
			msg:      "req: second",
		},
		{
			name:   "unprefixed line",
			reply:  "Invalid dispatcher",
			errors: []string{"Invalid dispatcher"},
			msg:    "req: Invalid dispatcher",
		},
		{name: "ok inside a longer reply", reply: "ok, but not really", errors: []string{"ok, but not really"}, msg: "req: ok, but not really"},
		{name: "empty", reply: "", msg: "req: unexpected reply"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := parseReply("req", tt.reply)
			if tt.ok {
				if err != nil {
					t.Fatalf("parseReply = %v, want nil", err)
				}
				return
			}

			var re *ReplyError
			if !errors.As(err, &re) {
				t.Fatalf("parseReply = %T, want *ReplyError", err)
			}
			if re.Request != "req" {
				t.Errorf("Request = %q, want req", re.Request)
			}
			if !equalStrings(re.Errors, tt.errors) {
				t.Errorf("Errors = %q, want %q", re.Errors, tt.errors)
			}
			if !equalStrings(re.Warnings, tt.warnings) {
				t.Errorf("Warnings = %q, want %q", re.Warnings, tt.warnings)
			}
			if re.Error() != tt.msg {
				t.Errorf("Error() = %q, want %q", re.Error(), tt.msg)
			}
		})
	}
}

func TestDispatch(t *testing.T) {
	t.Parallel()

	t.Run("prefixes the dispatch keyword", func(t *testing.T) {
		t.Parallel()

		rec := &commandRecorder{reply: "ok"}
		opts := commandOpts(t, rec)

		if err := Dispatch(context.Background(), FocusWorkspace{Workspace: WorkspaceByID(3)}, opts...); err != nil {
			t.Fatalf("Dispatch: %v", err)
		}

		want := `/dispatch hl.dsp.focus({ workspace = "3" })`
		if got := rec.commands(); len(got) != 1 || got[0] != want {
			t.Errorf("server received %q, want [%q]", got, want)
		}
	})

	t.Run("Request sends the request unchanged", func(t *testing.T) {
		t.Parallel()

		rec := &commandRecorder{reply: `{"class":"foot"}`}
		opts := commandOpts(t, rec)

		reply, err := Request(context.Background(), "j/activewindow", opts...)
		if err != nil {
			t.Fatalf("Request: %v", err)
		}
		if string(reply) != `{"class":"foot"}` {
			t.Errorf("reply = %q, want the raw JSON", reply)
		}
		if got := rec.commands(); len(got) != 1 || got[0] != "j/activewindow" {
			t.Errorf("server received %q, want [j/activewindow]", got)
		}
	})

	t.Run("reply other than ok is a ReplyError", func(t *testing.T) {
		t.Parallel()

		opts := commandOpts(t, &commandRecorder{reply: "Invalid dispatcher"})

		err := Dispatch(context.Background(), RawCommand("nope"), opts...)

		var re *ReplyError
		if !errors.As(err, &re) {
			t.Fatalf("Dispatch = %v, want *ReplyError", err)
		}
		if re.Request != "/dispatch nope" {
			t.Errorf("Request = %q, want /dispatch nope", re.Request)
		}
		if len(re.Errors) != 1 || re.Errors[0] != "Invalid dispatcher" {
			t.Errorf("Errors = %q, want [Invalid dispatcher]", re.Errors)
		}
	})

	t.Run("dials the command socket, not the event socket", func(t *testing.T) {
		t.Parallel()

		client, server := net.Pipe()
		defer server.Close()

		go func() {
			buf := make([]byte, 64)
			_, _ = server.Read(buf)
			_, _ = server.Write([]byte("ok"))
			server.Close()
		}()

		d := &fakeDialer{dial: func(context.Context) (net.Conn, error) { return client, nil }}
		err := Dispatch(context.Background(), RawCommand("nop"),
			WithDialer(d),
			WithSocketDir("/sockets"),
			WithCommandTimeout(time.Second),
			WithLogger(discardLogger()),
		)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}

		calls := d.dialed()
		if len(calls) != 1 {
			t.Fatalf("dialed %d times, want 1", len(calls))
		}
		if want := "/sockets/" + commandSocketName; calls[0].address != want {
			t.Errorf("command dial address = %q, want %q", calls[0].address, want)
		}
	})
}

func TestBatch(t *testing.T) {
	t.Parallel()

	cmds := []Command{RawCommand("hl.dsp.a()"), RawCommand("hl.dsp.b()"), RawCommand("hl.dsp.c()")}

	t.Run("one request for all commands", func(t *testing.T) {
		t.Parallel()

		rec := &commandRecorder{reply: "ok\n\n\nok\n\n\nok"}
		opts := commandOpts(t, rec)

		if err := Batch(context.Background(), cmds, opts...); err != nil {
			t.Fatalf("Batch: %v", err)
		}

		want := "[[BATCH]]/dispatch hl.dsp.a();/dispatch hl.dsp.b();/dispatch hl.dsp.c()"
		if got := rec.commands(); len(got) != 1 || got[0] != want {
			t.Errorf("server received %q, want [%q]", got, want)
		}
	})

	t.Run("maps each segment to its command", func(t *testing.T) {
		t.Parallel()

		opts := commandOpts(t, &commandRecorder{reply: "error: first\n\n\nok\n\n\nwarning: third"})

		err := Batch(context.Background(), cmds, opts...)
		if err == nil {
			t.Fatal("Batch succeeded, want two ReplyErrors")
		}

		var re *ReplyError
		if !errors.As(err, &re) {
			t.Fatalf("Batch = %T, want *ReplyError inside", err)
		}

		want := []string{"/dispatch hl.dsp.a(): first", "/dispatch hl.dsp.c(): third"}
		if got := strings.Split(err.Error(), "\n"); !equalStrings(got, want) {
			t.Errorf("Error() lines = %q, want %q", got, want)
		}
	})

	t.Run("segment count mismatch fails the whole batch", func(t *testing.T) {
		t.Parallel()

		opts := commandOpts(t, &commandRecorder{reply: "error: bad batch"})

		err := Batch(context.Background(), cmds, opts...)

		var re *ReplyError
		if !errors.As(err, &re) {
			t.Fatalf("Batch = %v, want *ReplyError", err)
		}
		if !strings.HasPrefix(re.Request, "[[BATCH]]") {
			t.Errorf("Request = %q, want the batch request", re.Request)
		}
		if len(re.Errors) != 1 || re.Errors[0] != "bad batch" {
			t.Errorf("Errors = %q, want [bad batch]", re.Errors)
		}
	})

	t.Run("no commands, no dial", func(t *testing.T) {
		t.Parallel()

		d := &fakeDialer{}
		if err := Batch(context.Background(), nil, WithDialer(d), WithSocketDir("/sockets")); err != nil {
			t.Fatalf("Batch: %v", err)
		}
		if calls := d.dialed(); len(calls) != 0 {
			t.Errorf("dialed %d times, want 0", len(calls))
		}
	})
}

func TestRequestMissingSignature(t *testing.T) {
	t.Setenv(hyprlandInstanceSignature, "")

	if _, err := Request(context.Background(), "nop"); !errors.Is(err, ErrNoHyprlandSocket) {
		t.Errorf("Request error = %v, want ErrNoHyprlandSocket", err)
	}
	if err := Dispatch(context.Background(), RawCommand("nop")); !errors.Is(err, ErrNoHyprlandSocket) {
		t.Errorf("Dispatch error = %v, want ErrNoHyprlandSocket", err)
	}
	if err := Batch(context.Background(), []Command{RawCommand("nop")}); !errors.Is(err, ErrNoHyprlandSocket) {
		t.Errorf("Batch error = %v, want ErrNoHyprlandSocket", err)
	}
}

func TestRequestErrors(t *testing.T) {
	t.Parallel()

	t.Run("nothing listening on the command socket", func(t *testing.T) {
		t.Parallel()

		opts := []Option{
			WithSocketDir(t.TempDir()),
			WithDialTimeout(time.Second),
			WithLogger(discardLogger()),
		}
		if _, err := Request(context.Background(), "nop", opts...); err == nil {
			t.Error("Request succeeded with no command socket")
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		t.Parallel()

		opts := commandOpts(t, &commandRecorder{reply: "ok"})

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if _, err := Request(ctx, "nop", opts...); err == nil {
			t.Error("Request succeeded with a cancelled context")
		}
	})

	t.Run("server never replies", func(t *testing.T) {
		t.Parallel()

		block := make(chan struct{})
		t.Cleanup(func() { close(block) })

		dir := socketDirWith(t, map[string]func(net.Conn){
			commandSocketName: func(conn net.Conn) {
				buf := make([]byte, 64)
				_, _ = conn.Read(buf)
				<-block // hold the connection open without replying
			},
		})

		start := time.Now()
		_, err := Request(context.Background(), "nop",
			WithSocketDir(dir),
			WithCommandTimeout(200*time.Millisecond),
			WithLogger(discardLogger()),
		)
		if err == nil {
			t.Error("Request succeeded against a silent server")
		}
		if elapsed := time.Since(start); elapsed > 3*time.Second {
			t.Errorf("Request took %v to give up, want the command timeout", elapsed)
		}
	})
}

// commandRecorder answers one command per connection and remembers what it read.
type commandRecorder struct {
	reply string

	mu       sync.Mutex
	received []string
}

func (r *commandRecorder) handle(conn net.Conn) {
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if n > 0 {
		r.mu.Lock()
		r.received = append(r.received, string(buf[:n]))
		r.mu.Unlock()
	}
	if n == 0 && err != nil {
		return
	}

	_, _ = conn.Write([]byte(r.reply))
}

func (r *commandRecorder) commands() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]string(nil), r.received...)
}

// commandOpts serves rec on a real command socket and returns the options that
// point the command functions at it.
func commandOpts(t *testing.T, rec *commandRecorder) []Option {
	t.Helper()

	dir := socketDirWith(t, map[string]func(net.Conn){
		commandSocketName: rec.handle,
	})

	return []Option{
		WithSocketDir(dir),
		WithDialTimeout(2 * time.Second),
		WithCommandTimeout(2 * time.Second),
		WithLogger(discardLogger()),
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
