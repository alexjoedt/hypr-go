package hypr

import "testing"

var workspaceCommandTests = []commandTest{
	{WorkspaceRename{Workspace: WorkspaceByID(1), Name: "web"}, `hl.dsp.workspace.rename({ workspace = "1", name = "web" })`},
	{WorkspaceRename{Workspace: WorkspaceByName("web")}, `hl.dsp.workspace.rename({ workspace = "name:web" })`},

	{WorkspaceChangeID{Workspace: WorkspaceByName("web"), ID: 7}, `hl.dsp.workspace.change_id({ workspace = "name:web", id = 7 })`},
	{WorkspaceChangeID{Workspace: WorkspaceByID(2), ID: 0}, `hl.dsp.workspace.change_id({ workspace = "2", id = 0 })`},

	{WorkspaceMove{Monitor: MonitorByName("DP-1")}, `hl.dsp.workspace.move({ monitor = "DP-1" })`},
	{WorkspaceMove{Workspace: WorkspaceByID(3), Monitor: MonitorByName("DP-1")}, `hl.dsp.workspace.move({ workspace = "3", monitor = "DP-1" })`},

	{WorkspaceSwapMonitors{Monitor1: MonitorByName("DP-1"), Monitor2: MonitorByName("HDMI-A-1")}, `hl.dsp.workspace.swap_monitors({ monitor1 = "DP-1", monitor2 = "HDMI-A-1" })`},

	{WorkspaceToggleSpecial(""), `hl.dsp.workspace.toggle_special()`},
	{WorkspaceToggleSpecial("magic"), `hl.dsp.workspace.toggle_special("magic")`},
}

func TestWorkspaceDispatchers(t *testing.T) {
	t.Parallel()

	runCommandTests(t, workspaceCommandTests)
}

var focusCommandTests = []commandTest{
	{Focus{Direction: DirectionLeft}, `hl.dsp.focus({ direction = "left" })`},
	{Focus{Direction: DirectionUp}, `hl.dsp.focus({ direction = "up" })`},

	{FocusMonitor{Monitor: MonitorByName("DP-1")}, `hl.dsp.focus({ monitor = "DP-1" })`},
	{FocusMonitor{Monitor: MonitorRelative(1)}, `hl.dsp.focus({ monitor = "+1" })`},

	{FocusWorkspace{Workspace: WorkspaceByID(3)}, `hl.dsp.focus({ workspace = "3" })`},
	{FocusWorkspace{Workspace: WorkspaceByID(0)}, `hl.dsp.focus({ workspace = "0" })`},
	{FocusWorkspace{Workspace: WorkspaceByID(-98)}, `hl.dsp.focus({ workspace = "-98" })`},
	{FocusWorkspace{Workspace: WorkspaceBySpecial("magic"), OnCurrentMonitor: true}, `hl.dsp.focus({ workspace = "special:magic", on_current_monitor = true })`},

	{FocusWindow{Window: WindowByClass("^kitty$")}, `hl.dsp.focus({ window = "class:^kitty$" })`},
	{FocusWindow{Window: WindowByAddress("5a1b")}, `hl.dsp.focus({ window = "address:0x5a1b" })`},

	{FocusUrgentOrLast(), `hl.dsp.focus({ urgent_or_last = true })`},
	{FocusLast(), `hl.dsp.focus({ last = true })`},
}

func TestFocusDispatchers(t *testing.T) {
	t.Parallel()

	runCommandTests(t, focusCommandTests)
}
