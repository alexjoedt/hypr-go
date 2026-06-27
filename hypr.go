package hypr

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	hyprlandInstanceSignature = "HYPRLAND_INSTANCE_SIGNATURE"
	xdgRuntimeDir             = "XDG_RUNTIME_DIR"

	eventSocketName   = ".socket2.sock"
	commandSocketName = ".socket.sock"

	defaultDialTimeout    = 5 * time.Second
	defaultReadTimeout    = 15 * time.Second
	defaultCommandTimeout = 3 * time.Second
)

// ErrNoHyprlandSocket is returned by Listen and the command functions when
// HYPRLAND_INSTANCE_SIGNATURE is empty and no WithSocketDir was given.
var ErrNoHyprlandSocket = errors.New("HYPRLAND_INSTANCE_SIGNATURE is empty")

// Dialer opens a connection to a Hyprland socket. *net.Dialer satisfies it, so
// the default needs no adapter; tests supply their own.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

var _ Dialer = (*net.Dialer)(nil)

// config is the resolved set of options for one call: defaults, then the
// options given, then the environment for the socket dir.
type config struct {
	dialer         Dialer
	socketDir      string
	dialTimeout    time.Duration
	readTimeout    time.Duration
	commandTimeout time.Duration
	log            *slog.Logger
}

func (c *config) eventSocket() string {
	return filepath.Join(c.socketDir, eventSocketName)
}

func (c *config) commandSocket() string {
	return filepath.Join(c.socketDir, commandSocketName)
}

// resolve applies opts over the defaults and fills in the socket dir from the
// environment when no option set it.
func resolve(opts []Option) (config, error) {
	c := config{
		dialer:         &net.Dialer{},
		dialTimeout:    defaultDialTimeout,
		readTimeout:    defaultReadTimeout,
		commandTimeout: defaultCommandTimeout,
		log:            slog.Default(),
	}
	for _, opt := range opts {
		opt(&c)
	}

	if c.socketDir == "" {
		dir, err := socketDir()
		if err != nil {
			return config{}, err
		}
		c.socketDir = dir
	}

	return c, nil
}

// Option configures Listen and the command functions. Options travel with the
// call and are resolved on each one.
type Option func(*config)

// WithDialer replaces the dialer used by Listen and the command functions.
// The default is a zero net.Dialer.
func WithDialer(d Dialer) Option {
	return func(o *config) { o.dialer = d }
}

// WithSocketDir overrides the directory holding the two Hyprland sockets.
// The default is $XDG_RUNTIME_DIR/hypr/$HYPRLAND_INSTANCE_SIGNATURE, read on
// each call.
func WithSocketDir(dir string) Option {
	return func(o *config) { o.socketDir = dir }
}

// WithDialTimeout bounds how long Listen and the command functions may take to
// connect.
func WithDialTimeout(d time.Duration) Option {
	return func(o *config) { o.dialTimeout = d }
}

// WithReadTimeout bounds a single read from the event socket, so it applies to
// Listen only. Hyprland can stay silent for a long time, so a read that times
// out is not an error: the loop sets a fresh deadline and reads again.
func WithReadTimeout(d time.Duration) Option {
	return func(o *config) { o.readTimeout = d }
}

// WithCommandTimeout bounds one round trip of Dispatch, Request or Batch on the
// command socket.
func WithCommandTimeout(d time.Duration) Option {
	return func(o *config) { o.commandTimeout = d }
}

// WithLogger replaces the logger used for debug output. The default is
// slog.Default().
func WithLogger(l *slog.Logger) Option {
	return func(o *config) { o.log = l }
}

// socketDir returns the directory Hyprland puts its sockets in for the running
// instance.
func socketDir() (string, error) {
	sig := os.Getenv(hyprlandInstanceSignature)
	if sig == "" {
		return "", ErrNoHyprlandSocket
	}

	return filepath.Join(os.Getenv(xdgRuntimeDir), "hypr", sig), nil
}
