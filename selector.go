package hypr

import "strconv"

// WindowSelector picks a window the way Hyprland's window selector grammar
// does. The zero value means the focused window wherever a dispatcher makes
// the window optional. The typed string is the escape hatch for grammar without
// a constructor ("negative:class:^foo$").
type WindowSelector string

// Window selectors that name a window by its state.
const (
	WindowActive   WindowSelector = "active"
	WindowFloating WindowSelector = "floating"
	WindowTiled    WindowSelector = "tiled"
)

// WindowByAddress selects the window with the given address.
func WindowByAddress(a WindowAddress) WindowSelector {
	return WindowSelector("address:0x" + a.String())
}

// WindowByClass selects the first window whose class matches the RE2 regex
// in full.
func WindowByClass(re string) WindowSelector {
	return WindowSelector("class:" + re)
}

// WindowByInitialClass selects by the class the window had when it opened.
func WindowByInitialClass(re string) WindowSelector {
	return WindowSelector("initialclass:" + re)
}

// WindowByTitle selects the first window whose title matches the RE2 regex
// in full.
func WindowByTitle(re string) WindowSelector {
	return WindowSelector("title:" + re)
}

// WindowByInitialTitle selects by the title the window had when it opened.
func WindowByInitialTitle(re string) WindowSelector {
	return WindowSelector("initialtitle:" + re)
}

// WindowByTag selects the first window carrying a tag that matches the RE2
// regex in full.
func WindowByTag(re string) WindowSelector {
	return WindowSelector("tag:" + re)
}

// WindowByPID selects the window owned by the process.
func WindowByPID(pid int) WindowSelector {
	return WindowSelector("pid:" + strconv.Itoa(pid))
}

// WorkspaceSelector names a workspace the way Hyprland's workspace grammar
// does. The typed string is the escape hatch for grammar without a constructor
// ("r+1", "emptynm").
type WorkspaceSelector string

// Workspace selectors relative to the active workspace.
const (
	WorkspacePrevious WorkspaceSelector = "previous"
	WorkspaceNext     WorkspaceSelector = "next"
	WorkspaceEmpty    WorkspaceSelector = "empty"
)

// WorkspaceByID selects a numbered workspace.
func WorkspaceByID(id int) WorkspaceSelector {
	return WorkspaceSelector(strconv.Itoa(id))
}

// WorkspaceByName selects a workspace by its name.
func WorkspaceByName(name string) WorkspaceSelector {
	return WorkspaceSelector("name:" + name)
}

// WorkspaceBySpecial selects a special workspace; the empty name is the unnamed
// one.
func WorkspaceBySpecial(name string) WorkspaceSelector {
	if name == "" {
		return "special"
	}

	return WorkspaceSelector("special:" + name)
}

// WorkspaceRelative selects the workspace n steps from the active one by ID,
// "+2" or "-1".
func WorkspaceRelative(n int) WorkspaceSelector {
	return WorkspaceSelector(signed(n))
}

// MonitorSelector names a monitor the way Hyprland's monitor grammar does. The
// typed string is the escape hatch for grammar without a constructor.
type MonitorSelector string

// MonitorCurrent selects the focused monitor.
const MonitorCurrent MonitorSelector = "current"

// MonitorByName selects a monitor by its output name, "DP-1".
func MonitorByName(name string) MonitorSelector {
	return MonitorSelector(name)
}

// MonitorByID selects a monitor by its ID.
func MonitorByID(id int) MonitorSelector {
	return MonitorSelector(strconv.Itoa(id))
}

// MonitorByDescription selects a monitor whose description starts with the
// prefix.
func MonitorByDescription(prefix string) MonitorSelector {
	return MonitorSelector("desc:" + prefix)
}

// MonitorRelative selects the monitor n places from the focused one in the
// monitor list, wrapping around.
func MonitorRelative(n int) MonitorSelector {
	return MonitorSelector(signed(n))
}

// MonitorInDirection selects the monitor in the given direction from the
// focused one.
func MonitorInDirection(d Direction) MonitorSelector {
	return MonitorSelector(d[:1])
}

// Direction is a cardinal direction for focus, move and swap dispatchers.
type Direction string

// The four directions.
const (
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
)

// signed formats n with an explicit sign, the form Hyprland's relative
// selectors need.
func signed(n int) string {
	if n >= 0 {
		return "+" + strconv.Itoa(n)
	}

	return strconv.Itoa(n)
}
