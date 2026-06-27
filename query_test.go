package hypr

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestQueryFixturesAreJSON catches a bad hand edit in testdata/query before a
// typed struct gets to see it.
func TestQueryFixturesAreJSON(t *testing.T) {
	t.Parallel()

	files, err := filepath.Glob("testdata/query/*.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 7 {
		t.Fatalf("found %d fixtures, want at least 7", len(files))
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		var v any
		if err := json.Unmarshal(data, &v); err != nil {
			t.Errorf("%s: %v", file, err)
		}
	}
}

// fixtureOpts serves one fixture file on a command socket and returns the
// options that point the query functions at it.
func fixtureOpts(t *testing.T, name string) []Option {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("testdata", "query", name))
	if err != nil {
		t.Fatal(err)
	}

	return commandOpts(t, &commandRecorder{reply: string(data)})
}

func TestQueryErrors(t *testing.T) {
	t.Parallel()

	t.Run("not JSON", func(t *testing.T) {
		t.Parallel()

		opts := commandOpts(t, &commandRecorder{reply: "unknown request\nsecond line"})

		_, err := Clients(context.Background(), opts...)
		if err == nil {
			t.Fatal("Clients succeeded on a non-JSON reply")
		}
		if !strings.Contains(err.Error(), `"unknown request"`) || strings.Contains(err.Error(), "second line") {
			t.Errorf("error = %q, want the first reply line only", err)
		}
		if !strings.HasPrefix(err.Error(), "j/clients: ") {
			t.Errorf("error = %q, want the request as prefix", err)
		}
	})

	t.Run("transport error", func(t *testing.T) {
		t.Parallel()

		_, err := Clients(context.Background(),
			WithSocketDir(t.TempDir()),
			WithDialTimeout(time.Second),
			WithLogger(discardLogger()),
		)
		if err == nil {
			t.Error("Clients succeeded with no socket")
		}
	})
}
