package hypr

import "syscall"

// FullscreenMode is the mode argument of WindowFullscreen. The zero value is
// Hyprland's default, fullscreen.
type FullscreenMode int

// The FullscreenMode values, in Hyprland's order.
const (
	FullscreenModeFull FullscreenMode = iota
	FullscreenModeMaximized
)

// String returns the Lua value: "fullscreen" or "maximized".
func (m FullscreenMode) String() string {
	if m == FullscreenModeMaximized {
		return "maximized"
	}

	return "fullscreen"
}

// FullscreenAction is the action argument of WindowFullscreen and
// WindowFullscreenState. The zero value is toggle.
type FullscreenAction int

// The FullscreenAction values, in Hyprland's order.
const (
	FullscreenToggle FullscreenAction = iota
	FullscreenSet
	FullscreenUnset
)

// String returns the Lua value: "toggle", "set" or "unset".
func (a FullscreenAction) String() string {
	switch a {
	case FullscreenSet:
		return "set"
	case FullscreenUnset:
		return "unset"
	default:
		return "toggle"
	}
}

// FullscreenState is one side of WindowFullscreenState: what the window
// reports to its client and what Hyprland renders.
type FullscreenState int

// The FullscreenState values, numbered the way Hyprland does.
const (
	FullscreenStateKeep      FullscreenState = -1
	FullscreenStateNone      FullscreenState = 0
	FullscreenStateMaximized FullscreenState = 1
	FullscreenStateFull      FullscreenState = 2
)

// ZOrder is the mode argument of WindowAlterZOrder.
type ZOrder int

// The ZOrder values.
const (
	ZOrderTop ZOrder = iota
	ZOrderBottom
)

// String returns the Lua value: "top" or "bottom".
func (z ZOrder) String() string {
	if z == ZOrderBottom {
		return "bottom"
	}

	return "top"
}

// WindowClose wraps hl.dsp.window.close: it asks the window to close. Window
// defaults to the focused one.
type WindowClose struct {
	Window WindowSelector
}

// Command implements Command.
func (c WindowClose) Command() string {
	return dispatch("window.close", optStr("window", string(c.Window)))
}

// WindowKill wraps hl.dsp.window.kill: it sends SIGKILL to the window's
// process. Window defaults to the focused one.
type WindowKill struct {
	Window WindowSelector
}

// Command implements Command.
func (c WindowKill) Command() string {
	return dispatch("window.kill", optStr("window", string(c.Window)))
}

// WindowSignal wraps hl.dsp.window.signal: it sends Signal, 1 to 31, to the
// window's process. Window defaults to the focused one.
type WindowSignal struct {
	Signal syscall.Signal
	Window WindowSelector
}

// Command implements Command.
func (c WindowSignal) Command() string {
	return dispatch("window.signal",
		num("signal", int(c.Signal)),
		optStr("window", string(c.Window)),
	)
}

// WindowFloat wraps hl.dsp.window.float: it toggles, enables or disables
// floating. Window defaults to the focused one.
type WindowFloat struct {
	Action Action
	Window WindowSelector
}

// Command implements Command.
func (c WindowFloat) Command() string {
	return dispatch("window.float",
		optAction("action", c.Action),
		optStr("window", string(c.Window)),
	)
}

// WindowFullscreen wraps hl.dsp.window.fullscreen. Mode defaults to
// fullscreen, Action to toggle; IgnoreLayout sends layout_aware = false.
// Window defaults to the focused one.
type WindowFullscreen struct {
	Mode         FullscreenMode
	Action       FullscreenAction
	IgnoreLayout bool
	Window       WindowSelector
}

// Command implements Command.
func (c WindowFullscreen) Command() string {
	return dispatch("window.fullscreen",
		optEnum("mode", c.Mode.String(), c.Mode == FullscreenModeFull),
		optEnum("action", c.Action.String(), c.Action == FullscreenToggle),
		optFalse("layout_aware", c.IgnoreLayout),
		optStr("window", string(c.Window)),
	)
}

// WindowFullscreenState wraps hl.dsp.window.fullscreen_state: it sets the
// internal (rendered) and client (reported) fullscreen state separately.
// Action is always sent, so the zero value toggles; use FullscreenSet for
// Hyprland's own default. IgnoreLayout sends layout_aware = false. Window
// defaults to the focused one.
type WindowFullscreenState struct {
	Internal     FullscreenState
	Client       FullscreenState
	Action       FullscreenAction
	IgnoreLayout bool
	Window       WindowSelector
}

// Command implements Command.
func (c WindowFullscreenState) Command() string {
	return dispatch("window.fullscreen_state",
		num("internal", int(c.Internal)),
		num("client", int(c.Client)),
		str("action", c.Action.String()),
		optFalse("layout_aware", c.IgnoreLayout),
		optStr("window", string(c.Window)),
	)
}

// WindowPseudo wraps hl.dsp.window.pseudo: it toggles, enables or disables
// pseudotiling. Window defaults to the focused one.
type WindowPseudo struct {
	Action Action
	Window WindowSelector
}

// Command implements Command.
func (c WindowPseudo) Command() string {
	return dispatch("window.pseudo",
		optAction("action", c.Action),
		optStr("window", string(c.Window)),
	)
}

// WindowMove wraps hl.dsp.window.move({ direction = … }): it moves the window
// in a direction within the layout. GroupAware moves into a neighbouring group
// instead of past it. Window defaults to the focused one.
type WindowMove struct {
	Direction  Direction
	GroupAware bool
	Window     WindowSelector
}

// Command implements Command.
func (c WindowMove) Command() string {
	return dispatch("window.move",
		str("direction", string(c.Direction)),
		optBool("group_aware", c.GroupAware),
		optStr("window", string(c.Window)),
	)
}

// WindowMoveTo wraps hl.dsp.window.move({ x = …, y = … }): it moves a floating
// window to a position, or by an offset when Relative is set. Window defaults
// to the focused one.
type WindowMoveTo struct {
	X, Y     int
	Relative bool
	Window   WindowSelector
}

// Command implements Command.
func (c WindowMoveTo) Command() string {
	return dispatch("window.move",
		num("x", c.X),
		num("y", c.Y),
		optBool("relative", c.Relative),
		optStr("window", string(c.Window)),
	)
}

// WindowMoveToWorkspace wraps hl.dsp.window.move({ workspace = … }): it moves
// the window to a workspace and follows it unless NoFollow is set. Window
// defaults to the focused one.
type WindowMoveToWorkspace struct {
	Workspace WorkspaceSelector
	NoFollow  bool
	Window    WindowSelector
}

// Command implements Command.
func (c WindowMoveToWorkspace) Command() string {
	return dispatch("window.move",
		str("workspace", string(c.Workspace)),
		optFalse("follow", c.NoFollow),
		optStr("window", string(c.Window)),
	)
}

// WindowMoveToMonitor wraps hl.dsp.window.move({ monitor = … }): it moves the
// window to a monitor and follows it unless NoFollow is set. Window defaults
// to the focused one.
type WindowMoveToMonitor struct {
	Monitor  MonitorSelector
	NoFollow bool
	Window   WindowSelector
}

// Command implements Command.
func (c WindowMoveToMonitor) Command() string {
	return dispatch("window.move",
		str("monitor", string(c.Monitor)),
		optFalse("follow", c.NoFollow),
		optStr("window", string(c.Window)),
	)
}

// WindowMoveIntoGroup wraps hl.dsp.window.move({ into_group = … }): it moves
// the window into the group in that direction. Window defaults to the focused
// one.
type WindowMoveIntoGroup struct {
	Direction Direction
	Window    WindowSelector
}

// Command implements Command.
func (c WindowMoveIntoGroup) Command() string {
	return dispatch("window.move",
		str("into_group", string(c.Direction)),
		optStr("window", string(c.Window)),
	)
}

// WindowMoveIntoOrCreateGroup wraps hl.dsp.window.move({ into_or_create_group
// = … }): like WindowMoveIntoGroup, but creates a group with the neighbour
// when there is none. Window defaults to the focused one.
type WindowMoveIntoOrCreateGroup struct {
	Direction Direction
	Window    WindowSelector
}

// Command implements Command.
func (c WindowMoveIntoOrCreateGroup) Command() string {
	return dispatch("window.move",
		str("into_or_create_group", string(c.Direction)),
		optStr("window", string(c.Window)),
	)
}

// WindowMoveOutOfGroup wraps hl.dsp.window.move({ out_of_group = … }): it moves
// the window out of its group, in Direction when set. Window defaults to the
// focused one.
type WindowMoveOutOfGroup struct {
	Direction Direction
	Window    WindowSelector
}

// Command implements Command.
func (c WindowMoveOutOfGroup) Command() string {
	out := field{key: "out_of_group", value: "true"}
	if c.Direction != "" {
		out.value = quote(string(c.Direction))
	}

	return dispatch("window.move", out, optStr("window", string(c.Window)))
}

// WindowSwap wraps hl.dsp.window.swap({ direction = … }): it swaps the window
// with its neighbour in that direction. Window defaults to the focused one.
type WindowSwap struct {
	Direction Direction
	Window    WindowSelector
}

// Command implements Command.
func (c WindowSwap) Command() string {
	return dispatch("window.swap",
		str("direction", string(c.Direction)),
		optStr("window", string(c.Window)),
	)
}

// WindowSwapWith wraps hl.dsp.window.swap({ target = … }): it swaps the window
// with Target. Window defaults to the focused one.
type WindowSwapWith struct {
	Target WindowSelector
	Window WindowSelector
}

// Command implements Command.
func (c WindowSwapWith) Command() string {
	return dispatch("window.swap",
		str("target", string(c.Target)),
		optStr("window", string(c.Window)),
	)
}

// WindowSwapNext wraps hl.dsp.window.swap({ next = true }): it swaps the
// focused window with the next one in the layout.
func WindowSwapNext() Command {
	return RawCommand(dispatch("window.swap", boolean("next", true)))
}

// WindowSwapPrev wraps hl.dsp.window.swap({ prev = true }): it swaps the
// focused window with the previous one in the layout.
func WindowSwapPrev() Command {
	return RawCommand(dispatch("window.swap", boolean("prev", true)))
}

// WindowCenter wraps hl.dsp.window.center: it centers a floating window on
// its monitor. Window defaults to the focused one.
type WindowCenter struct {
	Window WindowSelector
}

// Command implements Command.
func (c WindowCenter) Command() string {
	return dispatch("window.center", optStr("window", string(c.Window)))
}

// WindowCycleNext wraps hl.dsp.window.cycle_next: it focuses the next window
// on the workspace, or the previous one with Backward. Tiled and Floating
// restrict the cycle to that kind. Window defaults to the focused one.
type WindowCycleNext struct {
	Backward bool
	Tiled    bool
	Floating bool
	Window   WindowSelector
}

// Command implements Command.
func (c WindowCycleNext) Command() string {
	return dispatch("window.cycle_next",
		optFalse("next", c.Backward),
		optBool("tiled", c.Tiled),
		optBool("floating", c.Floating),
		optStr("window", string(c.Window)),
	)
}

// WindowTag wraps hl.dsp.window.tag: "name" toggles the tag, "+name" sets it,
// "-name" unsets it. Window defaults to the focused one.
type WindowTag struct {
	Tag    string
	Window WindowSelector
}

// Command implements Command.
func (c WindowTag) Command() string {
	return dispatch("window.tag",
		str("tag", c.Tag),
		optStr("window", string(c.Window)),
	)
}

// WindowClearTags wraps hl.dsp.window.clear_tags: it removes every tag from
// the window. Window defaults to the focused one.
type WindowClearTags struct {
	Window WindowSelector
}

// Command implements Command.
func (c WindowClearTags) Command() string {
	return dispatch("window.clear_tags", optStr("window", string(c.Window)))
}

// WindowToggleSwallow wraps hl.dsp.window.toggle_swallow: it toggles the
// focused window's swallow state.
func WindowToggleSwallow() Command {
	return RawCommand(dispatch("window.toggle_swallow"))
}

// WindowPin wraps hl.dsp.window.pin: it toggles, enables or disables pinning
// of a floating window to every workspace. Window defaults to the focused one.
type WindowPin struct {
	Action Action
	Window WindowSelector
}

// Command implements Command.
func (c WindowPin) Command() string {
	return dispatch("window.pin",
		optAction("action", c.Action),
		optStr("window", string(c.Window)),
	)
}

// WindowBringToTop wraps hl.dsp.window.bring_to_top: it raises the focused
// floating window above the others.
func WindowBringToTop() Command {
	return RawCommand(dispatch("window.bring_to_top"))
}

// WindowAlterZOrder wraps hl.dsp.window.alter_zorder: it sends a floating
// window to the top or the bottom of the stack. Window defaults to the
// focused one.
type WindowAlterZOrder struct {
	Mode   ZOrder
	Window WindowSelector
}

// Command implements Command.
func (c WindowAlterZOrder) Command() string {
	return dispatch("window.alter_zorder",
		str("mode", c.Mode.String()),
		optStr("window", string(c.Window)),
	)
}

// WindowSetProp wraps hl.dsp.window.set_prop: it applies a dynamic window
// rule effect such as "no_anim" or "opacity" with Value. Window defaults to
// the focused one.
type WindowSetProp struct {
	Prop   string
	Value  string
	Window WindowSelector
}

// Command implements Command.
func (c WindowSetProp) Command() string {
	return dispatch("window.set_prop",
		str("prop", c.Prop),
		str("value", c.Value),
		optStr("window", string(c.Window)),
	)
}

// WindowDenyFromGroup wraps hl.dsp.window.deny_from_group: it toggles,
// enables or disables the focused window's refusal to be grouped.
type WindowDenyFromGroup struct {
	Action Action
}

// Command implements Command.
func (c WindowDenyFromGroup) Command() string {
	return dispatch("window.deny_from_group", optAction("action", c.Action))
}

// WindowDrag wraps hl.dsp.window.drag: it starts an interactive mouse drag,
// for mouse binds.
func WindowDrag() Command {
	return RawCommand(dispatch("window.drag"))
}

// WindowResize wraps hl.dsp.window.resize without a size: it starts an
// interactive mouse resize, for mouse binds. KeepAspectRatio locks the
// window's aspect ratio during the drag.
type WindowResize struct {
	KeepAspectRatio bool
}

// Command implements Command.
func (c WindowResize) Command() string {
	return dispatch("window.resize", optBool("keep_aspect_ratio", c.KeepAspectRatio))
}

// WindowResizeTo wraps hl.dsp.window.resize({ x = …, y = … }): it resizes the
// window to a size, or by an offset when Relative is set. Window defaults to
// the focused one.
type WindowResizeTo struct {
	X, Y     int
	Relative bool
	Window   WindowSelector
}

// Command implements Command.
func (c WindowResizeTo) Command() string {
	return dispatch("window.resize",
		num("x", c.X),
		num("y", c.Y),
		optBool("relative", c.Relative),
		optStr("window", string(c.Window)),
	)
}
