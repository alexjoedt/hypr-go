package hypr

import (
	"errors"
	"testing"
	"time"
)

func TestSocketDir(t *testing.T) {
	t.Run("signature set", func(t *testing.T) {
		t.Setenv(xdgRuntimeDir, "/run/user/1000")
		t.Setenv(hyprlandInstanceSignature, "abc123")

		got, err := socketDir()
		if err != nil {
			t.Fatalf("socketDir: %v", err)
		}
		if want := "/run/user/1000/hypr/abc123"; got != want {
			t.Errorf("socketDir() = %q, want %q", got, want)
		}
	})

	t.Run("signature empty", func(t *testing.T) {
		t.Setenv(xdgRuntimeDir, "/run/user/1000")
		t.Setenv(hyprlandInstanceSignature, "")

		if _, err := socketDir(); !errors.Is(err, ErrNoHyprlandSocket) {
			t.Errorf("socketDir() error = %v, want ErrNoHyprlandSocket", err)
		}
	})
}

func TestResolve(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		t.Setenv(xdgRuntimeDir, "/run/user/1000")
		t.Setenv(hyprlandInstanceSignature, "abc123")

		cfg, err := resolve(nil)
		if err != nil {
			t.Fatalf("resolve: %v", err)
		}
		if cfg.eventSocket() != "/run/user/1000/hypr/abc123/"+eventSocketName {
			t.Errorf("eventSocket() = %q", cfg.eventSocket())
		}
		if cfg.commandSocket() != "/run/user/1000/hypr/abc123/"+commandSocketName {
			t.Errorf("commandSocket() = %q", cfg.commandSocket())
		}
		if cfg.dialTimeout != defaultDialTimeout || cfg.readTimeout != defaultReadTimeout || cfg.commandTimeout != defaultCommandTimeout {
			t.Errorf("timeouts = %v %v %v, want the defaults", cfg.dialTimeout, cfg.readTimeout, cfg.commandTimeout)
		}
		if cfg.dialer == nil || cfg.log == nil {
			t.Error("dialer or logger is nil")
		}
	})

	t.Run("options win over the environment", func(t *testing.T) {
		t.Setenv(hyprlandInstanceSignature, "")

		cfg, err := resolve([]Option{
			WithSocketDir("/sockets"),
			WithDialTimeout(time.Second),
			WithReadTimeout(2 * time.Second),
			WithCommandTimeout(3 * time.Second),
			WithDialer(&fakeDialer{}),
			WithLogger(discardLogger()),
		})
		if err != nil {
			t.Fatalf("resolve: %v", err)
		}
		if cfg.socketDir != "/sockets" {
			t.Errorf("socketDir = %q, want /sockets", cfg.socketDir)
		}
		if cfg.dialTimeout != time.Second || cfg.readTimeout != 2*time.Second || cfg.commandTimeout != 3*time.Second {
			t.Errorf("timeouts = %v %v %v", cfg.dialTimeout, cfg.readTimeout, cfg.commandTimeout)
		}
		if _, ok := cfg.dialer.(*fakeDialer); !ok {
			t.Errorf("dialer = %T, want *fakeDialer", cfg.dialer)
		}
	})

	t.Run("missing signature", func(t *testing.T) {
		t.Setenv(hyprlandInstanceSignature, "")

		if _, err := resolve(nil); !errors.Is(err, ErrNoHyprlandSocket) {
			t.Errorf("resolve error = %v, want ErrNoHyprlandSocket", err)
		}
	})
}
