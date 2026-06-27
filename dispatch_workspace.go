package hypr

// WorkspaceRename wraps hl.dsp.workspace.rename: it sets the display name of
// an existing workspace; an empty Name clears it.
type WorkspaceRename struct {
	Workspace WorkspaceSelector
	Name      string
}

// Command implements Command.
func (c WorkspaceRename) Command() string {
	return dispatch("workspace.rename",
		str("workspace", string(c.Workspace)),
		optStr("name", c.Name),
	)
}

// WorkspaceChangeID wraps hl.dsp.workspace.change_id: it gives an existing
// workspace a new, unused ID.
type WorkspaceChangeID struct {
	Workspace WorkspaceSelector
	ID        int
}

// Command implements Command.
func (c WorkspaceChangeID) Command() string {
	return dispatch("workspace.change_id",
		str("workspace", string(c.Workspace)),
		num("id", c.ID),
	)
}

// WorkspaceMove wraps hl.dsp.workspace.move: it moves a workspace to a
// monitor. Workspace defaults to the focused monitor's current one.
type WorkspaceMove struct {
	Monitor   MonitorSelector
	Workspace WorkspaceSelector
}

// Command implements Command.
func (c WorkspaceMove) Command() string {
	return dispatch("workspace.move",
		optStr("workspace", string(c.Workspace)),
		str("monitor", string(c.Monitor)),
	)
}

// WorkspaceSwapMonitors wraps hl.dsp.workspace.swap_monitors: it swaps the
// active workspaces of two monitors.
type WorkspaceSwapMonitors struct {
	Monitor1 MonitorSelector
	Monitor2 MonitorSelector
}

// Command implements Command.
func (c WorkspaceSwapMonitors) Command() string {
	return dispatch("workspace.swap_monitors",
		str("monitor1", string(c.Monitor1)),
		str("monitor2", string(c.Monitor2)),
	)
}

// WorkspaceToggleSpecial wraps hl.dsp.workspace.toggle_special: it shows or
// hides the special workspace with that name; the empty name is the unnamed
// one.
func WorkspaceToggleSpecial(name string) Command {
	if name == "" {
		return RawCommand(positional("workspace.toggle_special"))
	}

	return RawCommand(positional("workspace.toggle_special", quote(name)))
}

// Focus wraps hl.dsp.focus({ direction = … }): it focuses the window in that
// direction.
type Focus struct {
	Direction Direction
}

// Command implements Command.
func (c Focus) Command() string {
	return dispatch("focus", str("direction", string(c.Direction)))
}

// FocusMonitor wraps hl.dsp.focus({ monitor = … }): it focuses the monitor.
type FocusMonitor struct {
	Monitor MonitorSelector
}

// Command implements Command.
func (c FocusMonitor) Command() string {
	return dispatch("focus", str("monitor", string(c.Monitor)))
}

// FocusWorkspace wraps hl.dsp.focus({ workspace = … }): it switches to the
// workspace, creating it when needed. OnCurrentMonitor keeps the workspace on
// the focused monitor instead of following it.
type FocusWorkspace struct {
	Workspace        WorkspaceSelector
	OnCurrentMonitor bool
}

// Command implements Command.
func (c FocusWorkspace) Command() string {
	return dispatch("focus",
		str("workspace", string(c.Workspace)),
		optBool("on_current_monitor", c.OnCurrentMonitor),
	)
}

// FocusWindow wraps hl.dsp.focus({ window = … }): it focuses the window,
// switching workspace when needed.
type FocusWindow struct {
	Window WindowSelector
}

// Command implements Command.
func (c FocusWindow) Command() string {
	return dispatch("focus", str("window", string(c.Window)))
}

// FocusUrgentOrLast wraps hl.dsp.focus({ urgent_or_last = true }): it focuses
// the urgent window, else the last focused one.
func FocusUrgentOrLast() Command {
	return RawCommand(dispatch("focus", boolean("urgent_or_last", true)))
}

// FocusLast wraps hl.dsp.focus({ last = true }): it focuses the previously
// focused window.
func FocusLast() Command {
	return RawCommand(dispatch("focus", boolean("last", true)))
}
