package hypr

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"
)

// Command is a dispatcher expression, hl.dsp.…(...), as the /dispatch request
// expects it.
type Command interface {
	Command() string
}

// RawCommand implements Command verbatim. It is the escape hatch for
// dispatchers this version has no constructor for.
type RawCommand string

// Command returns the expression verbatim.
func (rc RawCommand) Command() string {
	return string(rc)
}

// ReplyError is a reply Hyprland did not answer with ok. Lines are split by
// their "error:" or "warning:" prefix; anything else lands in Errors. Warnings
// alone are still a failure: hl.dispatch reports ok = false for them.
type ReplyError struct {
	// Request is the request as it went over the socket.
	Request  string
	Errors   []string
	Warnings []string
}

// Error returns the request and its first error line, else its first warning.
func (e *ReplyError) Error() string {
	msg := "unexpected reply"
	switch {
	case len(e.Errors) > 0:
		msg = e.Errors[0]
	case len(e.Warnings) > 0:
		msg = e.Warnings[0]
	}

	return e.Request + ": " + msg
}

// batchPrefix starts a request that carries several commands; batchSeparator
// splits the commands, batchReplySeparator the reply segments. Verified against
// Hyprland 0.56.2 on 2026-09-11: three commands answer
// "error: …\n\n\nok\n\n\nerror: …".
const (
	batchPrefix         = "[[BATCH]]"
	batchSeparator      = ";"
	batchReplySeparator = "\n\n\n"
	dispatchPrefix      = "/dispatch "
)

// Dispatch runs cmd as a Hyprland dispatcher. A reply other than ok is a
// *ReplyError.
func Dispatch(ctx context.Context, cmd Command, opts ...Option) error {
	request := dispatchPrefix + cmd.Command()

	reply, err := Request(ctx, request, opts...)
	if err != nil {
		return err
	}

	return parseReply(request, string(reply))
}

// Request sends request to the command socket unchanged, so it can carry any
// hyprctl request ("j/clients", "/reload"), and returns the reply bytes.
func Request(ctx context.Context, request string, opts ...Option) ([]byte, error) {
	cfg, err := resolve(opts)
	if err != nil {
		return nil, err
	}

	return cfg.roundTrip(ctx, request)
}

// Batch runs cmds in one round trip. Every command that did not answer ok
// comes back as a *ReplyError; several are joined with errors.Join.
func Batch(ctx context.Context, cmds []Command, opts ...Option) error {
	if len(cmds) == 0 {
		return nil
	}

	requests := make([]string, len(cmds))
	for i, cmd := range cmds {
		requests[i] = dispatchPrefix + cmd.Command()
	}
	request := batchPrefix + strings.Join(requests, batchSeparator)

	reply, err := Request(ctx, request, opts...)
	if err != nil {
		return err
	}

	segments := strings.Split(string(reply), batchReplySeparator)
	if len(segments) != len(requests) {
		// Hyprland rejected the batch as a whole rather than command by command.
		return parseReply(request, string(reply))
	}

	errs := make([]error, 0, len(segments))
	for i, segment := range segments {
		if err := parseReply(requests[i], segment); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// parseReply turns one reply into nil for ok and a *ReplyError for anything
// else.
func parseReply(request, reply string) error {
	reply = strings.TrimSpace(reply)
	if strings.EqualFold(reply, "ok") {
		return nil
	}

	re := &ReplyError{Request: request}
	for line := range strings.SplitSeq(reply, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
		case strings.HasPrefix(line, "error:"):
			re.Errors = append(re.Errors, strings.TrimSpace(strings.TrimPrefix(line, "error:")))
		case strings.HasPrefix(line, "warning:"):
			re.Warnings = append(re.Warnings, strings.TrimSpace(strings.TrimPrefix(line, "warning:")))
		default:
			re.Errors = append(re.Errors, line)
		}
	}

	return re
}

// roundTrip opens a fresh connection to the command socket, sends request and
// reads the reply until Hyprland closes the connection: one request per dial.
func (c *config) roundTrip(ctx context.Context, request string) ([]byte, error) {
	dialCtx, cancel := context.WithTimeout(ctx, c.dialTimeout)
	defer cancel()

	conn, err := c.dialer.DialContext(dialCtx, "unix", c.commandSocket())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c.log.Debug("connected to socket", "socket", c.commandSocket())

	// One deadline for the whole round trip. Without it a server that neither
	// replies nor closes leaves io.ReadAll blocked forever.
	_ = conn.SetDeadline(time.Now().Add(c.commandTimeout))

	if _, err := conn.Write([]byte(request)); err != nil {
		return nil, err
	}
	c.log.Debug("request sent", "request", request)

	return io.ReadAll(conn)
}
