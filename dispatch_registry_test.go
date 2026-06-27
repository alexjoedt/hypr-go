package hypr

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

// TestDispatcherCoverage keeps the constructors in step with the dispatchers
// Hyprland documents, listed in testdata/dispatchers.txt. Coverage is derived
// from the wire-string tables of the family tests: a constructor counts once a
// row exercises it, so the package itself carries no registry.
func TestDispatcherCoverage(t *testing.T) {
	t.Parallel()

	covered := make(map[string]bool)
	for _, tt := range allCommandTests() {
		path := dispatcherPath(t, tt.cmd)
		covered[path] = true
		checkDispatcherName(t, tt.cmd, path)
	}

	documented := make(map[string]bool)
	for _, line := range readLines(t, "testdata/dispatchers.txt") {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		path := strings.TrimPrefix(name, "hl.dsp.")
		documented[path] = true

		if !covered[path] {
			t.Errorf("dispatcher %s has no constructor", name)
		}
	}

	for path := range covered {
		if !documented[path] {
			t.Errorf("dispatcher hl.dsp.%s has a constructor but is not in testdata/dispatchers.txt", path)
		}
	}

	if len(covered) != 52 {
		t.Errorf("%d dispatchers covered, want 52", len(covered))
	}
}

// allCommandTests joins the wire-string tables of every family file.
func allCommandTests() []commandTest {
	return slices.Concat(
		sessionCommandTests,
		workspaceCommandTests,
		focusCommandTests,
		windowCommandTests,
		groupCommandTests,
		cursorCommandTests,
	)
}

// dispatcherPath returns the Lua path of cmd, the part of its expression
// between hl.dsp. and the opening parenthesis.
func dispatcherPath(t *testing.T, cmd Command) string {
	t.Helper()

	expr := cmd.Command()
	rest, ok := strings.CutPrefix(expr, "hl.dsp.")
	if !ok {
		t.Fatalf("%T: expression %q does not start with hl.dsp.", cmd, expr)
	}

	path, _, ok := strings.Cut(rest, "(")
	if !ok {
		t.Fatalf("%T: expression %q has no argument list", cmd, expr)
	}

	return path
}

// checkDispatcherName enforces the rule in doc.go for struct dispatchers: the
// Go name starts with the Lua path in CamelCase, an overload mode may follow.
// Function constructors return a RawCommand and carry no Go name to check.
func checkDispatcherName(t *testing.T, cmd Command, path string) {
	t.Helper()

	typ := reflect.TypeOf(cmd)
	if typ.Kind() != reflect.Struct {
		return
	}

	want := strings.NewReplacer(".", "", "_", "").Replace(path)
	if !strings.HasPrefix(strings.ToLower(typ.Name()), want) {
		t.Errorf("dispatcher hl.dsp.%s is %s, want a name starting with %s", path, typ.Name(), want)
	}
}
