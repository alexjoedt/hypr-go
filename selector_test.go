package hypr

import "testing"

func TestWindowSelector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sel  WindowSelector
		want string
	}{
		{"zero value", WindowSelector(""), ""},
		{"active", WindowActive, "active"},
		{"floating", WindowFloating, "floating"},
		{"tiled", WindowTiled, "tiled"},
		{"address", WindowByAddress("56157f7c25a0"), "address:0x56157f7c25a0"},
		{"class", WindowByClass("^kitty$"), "class:^kitty$"},
		{"initial class", WindowByInitialClass("^kitty$"), "initialclass:^kitty$"},
		{"title", WindowByTitle(".*vim.*"), "title:.*vim.*"},
		{"initial title", WindowByInitialTitle("^zsh$"), "initialtitle:^zsh$"},
		{"tag", WindowByTag("urgent"), "tag:urgent"},
		{"pid", WindowByPID(4242), "pid:4242"},
		{"escape hatch", WindowSelector("negative:class:^foo$"), "negative:class:^foo$"},
	}

	for _, tt := range tests {
		if string(tt.sel) != tt.want {
			t.Errorf("%s: selector = %q, want %q", tt.name, tt.sel, tt.want)
		}
	}
}

func TestWorkspaceSelector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sel  WorkspaceSelector
		want string
	}{
		{"id", WorkspaceByID(3), "3"},
		{"negative id", WorkspaceByID(-98), "-98"},
		{"name", WorkspaceByName("web"), "name:web"},
		{"special unnamed", WorkspaceBySpecial(""), "special"},
		{"special named", WorkspaceBySpecial("magic"), "special:magic"},
		{"relative forward", WorkspaceRelative(2), "+2"},
		{"relative back", WorkspaceRelative(-1), "-1"},
		{"relative zero", WorkspaceRelative(0), "+0"},
		{"previous", WorkspacePrevious, "previous"},
		{"next", WorkspaceNext, "next"},
		{"empty", WorkspaceEmpty, "empty"},
		{"escape hatch", WorkspaceSelector("emptynm"), "emptynm"},
	}

	for _, tt := range tests {
		if string(tt.sel) != tt.want {
			t.Errorf("%s: selector = %q, want %q", tt.name, tt.sel, tt.want)
		}
	}
}

func TestMonitorSelector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sel  MonitorSelector
		want string
	}{
		{"current", MonitorCurrent, "current"},
		{"name", MonitorByName("DP-1"), "DP-1"},
		{"id", MonitorByID(1), "1"},
		{"description", MonitorByDescription("Dell"), "desc:Dell"},
		{"relative forward", MonitorRelative(1), "+1"},
		{"relative back", MonitorRelative(-1), "-1"},
		{"left", MonitorInDirection(DirectionLeft), "l"},
		{"right", MonitorInDirection(DirectionRight), "r"},
		{"up", MonitorInDirection(DirectionUp), "u"},
		{"down", MonitorInDirection(DirectionDown), "d"},
	}

	for _, tt := range tests {
		if string(tt.sel) != tt.want {
			t.Errorf("%s: selector = %q, want %q", tt.name, tt.sel, tt.want)
		}
	}
}
