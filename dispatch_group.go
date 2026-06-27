package hypr

// Corner is the corner argument of CursorMoveToCorner, numbered the way
// Hyprland does.
type Corner int

// The Corner values, in Hyprland's order.
const (
	CornerBottomLeft Corner = iota
	CornerBottomRight
	CornerTopRight
	CornerTopLeft
)

// GroupToggle wraps hl.dsp.group.toggle: it turns the window into a group or
// dissolves its group. Window defaults to the focused one.
type GroupToggle struct {
	Window WindowSelector
}

// Command implements Command.
func (c GroupToggle) Command() string {
	return dispatch("group.toggle", optStr("window", string(c.Window)))
}

// GroupNext wraps hl.dsp.group.next: it focuses the next window in the group.
// Window defaults to the focused one.
type GroupNext struct {
	Window WindowSelector
}

// Command implements Command.
func (c GroupNext) Command() string {
	return dispatch("group.next", optStr("window", string(c.Window)))
}

// GroupPrev wraps hl.dsp.group.prev: it focuses the previous window in the
// group. Window defaults to the focused one.
type GroupPrev struct {
	Window WindowSelector
}

// Command implements Command.
func (c GroupPrev) Command() string {
	return dispatch("group.prev", optStr("window", string(c.Window)))
}

// GroupActive wraps hl.dsp.group.active: it focuses the window at Index in
// the group, counted from 1; 0 or less selects the last one. Window defaults
// to the focused one.
type GroupActive struct {
	Index  int
	Window WindowSelector
}

// Command implements Command.
func (c GroupActive) Command() string {
	return dispatch("group.active",
		num("index", c.Index),
		optStr("window", string(c.Window)),
	)
}

// GroupMoveWindow wraps hl.dsp.group.move_window: it moves the focused window
// one place forward in its group, or back with Backward.
type GroupMoveWindow struct {
	Backward bool
}

// Command implements Command.
func (c GroupMoveWindow) Command() string {
	return dispatch("group.move_window", optFalse("forward", c.Backward))
}

// GroupLock wraps hl.dsp.group.lock: it toggles, enables or disables the lock
// on every group.
type GroupLock struct {
	Action Action
}

// Command implements Command.
func (c GroupLock) Command() string {
	return dispatch("group.lock", optAction("action", c.Action))
}

// GroupLockActive wraps hl.dsp.group.lock_active: it toggles, enables or
// disables the lock on the focused group.
type GroupLockActive struct {
	Action Action
}

// Command implements Command.
func (c GroupLockActive) Command() string {
	return dispatch("group.lock_active", optAction("action", c.Action))
}

// CursorMoveToCorner wraps hl.dsp.cursor.move_to_corner: it warps the cursor
// to a corner of the window. Window defaults to the focused one.
type CursorMoveToCorner struct {
	Corner Corner
	Window WindowSelector
}

// Command implements Command.
func (c CursorMoveToCorner) Command() string {
	return dispatch("cursor.move_to_corner",
		num("corner", int(c.Corner)),
		optStr("window", string(c.Window)),
	)
}

// CursorMove wraps hl.dsp.cursor.move: it warps the cursor to a position in
// global layout coordinates.
type CursorMove struct {
	X, Y int
}

// Command implements Command.
func (c CursorMove) Command() string {
	return dispatch("cursor.move", num("x", c.X), num("y", c.Y))
}
