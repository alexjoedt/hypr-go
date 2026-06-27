package hypr

import (
	"syscall"
	"testing"
)

// commandTest is one row of a wire-string table: the command and the exact
// expression it must produce.
type commandTest struct {
	cmd  Command
	want string
}

// kitty is the window selector the wire-string tables share.
var kitty = WindowByClass("^kitty$")

func runCommandTests(t *testing.T, tests []commandTest) {
	t.Helper()

	for _, tt := range tests {
		if got := tt.cmd.Command(); got != tt.want {
			t.Errorf("%T%+v\n got %s\nwant %s", tt.cmd, tt.cmd, got, tt.want)
		}
	}
}

var windowCommandTests = []commandTest{
	{WindowClose{}, `hl.dsp.window.close()`},
	{WindowClose{Window: kitty}, `hl.dsp.window.close({ window = "class:^kitty$" })`},

	{WindowKill{}, `hl.dsp.window.kill()`},
	{WindowKill{Window: WindowActive}, `hl.dsp.window.kill({ window = "active" })`},

	{WindowSignal{Signal: syscall.SIGTERM}, `hl.dsp.window.signal({ signal = 15 })`},
	{WindowSignal{Signal: syscall.SIGHUP, Window: kitty}, `hl.dsp.window.signal({ signal = 1, window = "class:^kitty$" })`},

	{WindowFloat{}, `hl.dsp.window.float()`},
	{WindowFloat{Action: ActionEnable}, `hl.dsp.window.float({ action = "on" })`},
	{WindowFloat{Action: ActionDisable, Window: kitty}, `hl.dsp.window.float({ action = "off", window = "class:^kitty$" })`},

	{WindowFullscreen{}, `hl.dsp.window.fullscreen()`},
	{WindowFullscreen{Mode: FullscreenModeMaximized, Action: FullscreenSet}, `hl.dsp.window.fullscreen({ mode = "maximized", action = "set" })`},
	{WindowFullscreen{Action: FullscreenUnset, IgnoreLayout: true, Window: kitty}, `hl.dsp.window.fullscreen({ action = "unset", layout_aware = false, window = "class:^kitty$" })`},

	{WindowFullscreenState{Internal: FullscreenStateFull, Client: FullscreenStateNone, Action: FullscreenSet}, `hl.dsp.window.fullscreen_state({ internal = 2, client = 0, action = "set" })`},
	{WindowFullscreenState{Internal: FullscreenStateKeep, Client: FullscreenStateMaximized}, `hl.dsp.window.fullscreen_state({ internal = -1, client = 1, action = "toggle" })`},
	{WindowFullscreenState{Action: FullscreenUnset, IgnoreLayout: true, Window: kitty}, `hl.dsp.window.fullscreen_state({ internal = 0, client = 0, action = "unset", layout_aware = false, window = "class:^kitty$" })`},

	{WindowPseudo{}, `hl.dsp.window.pseudo()`},
	{WindowPseudo{Action: ActionEnable, Window: kitty}, `hl.dsp.window.pseudo({ action = "on", window = "class:^kitty$" })`},

	{WindowMove{Direction: DirectionLeft}, `hl.dsp.window.move({ direction = "left" })`},
	{WindowMove{Direction: DirectionDown, GroupAware: true, Window: kitty}, `hl.dsp.window.move({ direction = "down", group_aware = true, window = "class:^kitty$" })`},
	{WindowMoveTo{X: 100, Y: 200}, `hl.dsp.window.move({ x = 100, y = 200 })`},
	{WindowMoveTo{X: -10, Y: 0, Relative: true, Window: kitty}, `hl.dsp.window.move({ x = -10, y = 0, relative = true, window = "class:^kitty$" })`},
	{WindowMoveToWorkspace{Workspace: WorkspaceByID(3)}, `hl.dsp.window.move({ workspace = "3" })`},
	{WindowMoveToWorkspace{Workspace: WorkspaceBySpecial("magic"), NoFollow: true, Window: kitty}, `hl.dsp.window.move({ workspace = "special:magic", follow = false, window = "class:^kitty$" })`},
	{WindowMoveToMonitor{Monitor: MonitorByName("DP-1")}, `hl.dsp.window.move({ monitor = "DP-1" })`},
	{WindowMoveToMonitor{Monitor: MonitorInDirection(DirectionRight), NoFollow: true}, `hl.dsp.window.move({ monitor = "r", follow = false })`},
	{WindowMoveIntoGroup{Direction: DirectionLeft}, `hl.dsp.window.move({ into_group = "left" })`},
	{WindowMoveIntoGroup{Direction: DirectionUp, Window: kitty}, `hl.dsp.window.move({ into_group = "up", window = "class:^kitty$" })`},
	{WindowMoveIntoOrCreateGroup{Direction: DirectionRight}, `hl.dsp.window.move({ into_or_create_group = "right" })`},
	{WindowMoveOutOfGroup{}, `hl.dsp.window.move({ out_of_group = true })`},
	{WindowMoveOutOfGroup{Direction: DirectionDown, Window: kitty}, `hl.dsp.window.move({ out_of_group = "down", window = "class:^kitty$" })`},

	{WindowSwap{Direction: DirectionLeft}, `hl.dsp.window.swap({ direction = "left" })`},
	{WindowSwap{Direction: DirectionRight, Window: kitty}, `hl.dsp.window.swap({ direction = "right", window = "class:^kitty$" })`},
	{WindowSwapWith{Target: WindowByAddress("5a1b")}, `hl.dsp.window.swap({ target = "address:0x5a1b" })`},
	{WindowSwapWith{Target: kitty, Window: WindowActive}, `hl.dsp.window.swap({ target = "class:^kitty$", window = "active" })`},
	{WindowSwapNext(), `hl.dsp.window.swap({ next = true })`},
	{WindowSwapPrev(), `hl.dsp.window.swap({ prev = true })`},

	{WindowCenter{}, `hl.dsp.window.center()`},
	{WindowCenter{Window: kitty}, `hl.dsp.window.center({ window = "class:^kitty$" })`},

	{WindowCycleNext{}, `hl.dsp.window.cycle_next()`},
	{WindowCycleNext{Backward: true, Tiled: true}, `hl.dsp.window.cycle_next({ next = false, tiled = true })`},
	{WindowCycleNext{Floating: true, Window: kitty}, `hl.dsp.window.cycle_next({ floating = true, window = "class:^kitty$" })`},

	{WindowTag{Tag: "+urgent"}, `hl.dsp.window.tag({ tag = "+urgent" })`},
	{WindowTag{Tag: "", Window: kitty}, `hl.dsp.window.tag({ tag = "", window = "class:^kitty$" })`},

	{WindowClearTags{}, `hl.dsp.window.clear_tags()`},
	{WindowClearTags{Window: kitty}, `hl.dsp.window.clear_tags({ window = "class:^kitty$" })`},

	{WindowToggleSwallow(), `hl.dsp.window.toggle_swallow()`},

	{WindowPin{}, `hl.dsp.window.pin()`},
	{WindowPin{Action: ActionEnable, Window: kitty}, `hl.dsp.window.pin({ action = "on", window = "class:^kitty$" })`},

	{WindowBringToTop(), `hl.dsp.window.bring_to_top()`},

	{WindowAlterZOrder{}, `hl.dsp.window.alter_zorder({ mode = "top" })`},
	{WindowAlterZOrder{Mode: ZOrderBottom, Window: kitty}, `hl.dsp.window.alter_zorder({ mode = "bottom", window = "class:^kitty$" })`},

	{WindowSetProp{Prop: "no_anim", Value: "1"}, `hl.dsp.window.set_prop({ prop = "no_anim", value = "1" })`},
	{WindowSetProp{Prop: "opacity", Value: "0.5 override", Window: kitty}, `hl.dsp.window.set_prop({ prop = "opacity", value = "0.5 override", window = "class:^kitty$" })`},

	{WindowDenyFromGroup{}, `hl.dsp.window.deny_from_group()`},
	{WindowDenyFromGroup{Action: ActionEnable}, `hl.dsp.window.deny_from_group({ action = "on" })`},

	{WindowDrag(), `hl.dsp.window.drag()`},

	{WindowResize{}, `hl.dsp.window.resize()`},
	{WindowResize{KeepAspectRatio: true}, `hl.dsp.window.resize({ keep_aspect_ratio = true })`},
	{WindowResizeTo{X: 800, Y: 600}, `hl.dsp.window.resize({ x = 800, y = 600 })`},
	{WindowResizeTo{X: 100, Y: -50, Relative: true, Window: kitty}, `hl.dsp.window.resize({ x = 100, y = -50, relative = true, window = "class:^kitty$" })`},
}

func TestWindowDispatchers(t *testing.T) {
	t.Parallel()

	runCommandTests(t, windowCommandTests)
}

func TestWindowEnumStrings(t *testing.T) {
	t.Parallel()

	if FullscreenModeFull.String() != "fullscreen" || FullscreenModeMaximized.String() != "maximized" || FullscreenMode(9).String() != "fullscreen" {
		t.Error("FullscreenMode.String")
	}
	if FullscreenToggle.String() != "toggle" || FullscreenSet.String() != "set" || FullscreenUnset.String() != "unset" {
		t.Error("FullscreenAction.String")
	}
	if ZOrderTop.String() != "top" || ZOrderBottom.String() != "bottom" || ZOrder(9).String() != "top" {
		t.Error("ZOrder.String")
	}
}
