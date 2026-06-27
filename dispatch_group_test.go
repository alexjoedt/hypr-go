package hypr

import "testing"

var groupCommandTests = []commandTest{
	{GroupToggle{}, `hl.dsp.group.toggle()`},
	{GroupToggle{Window: kitty}, `hl.dsp.group.toggle({ window = "class:^kitty$" })`},

	{GroupNext{}, `hl.dsp.group.next()`},
	{GroupNext{Window: kitty}, `hl.dsp.group.next({ window = "class:^kitty$" })`},

	{GroupPrev{}, `hl.dsp.group.prev()`},
	{GroupPrev{Window: kitty}, `hl.dsp.group.prev({ window = "class:^kitty$" })`},

	{GroupActive{Index: 2}, `hl.dsp.group.active({ index = 2 })`},
	{GroupActive{Index: 0, Window: kitty}, `hl.dsp.group.active({ index = 0, window = "class:^kitty$" })`},

	{GroupMoveWindow{}, `hl.dsp.group.move_window()`},
	{GroupMoveWindow{Backward: true}, `hl.dsp.group.move_window({ forward = false })`},

	{GroupLock{}, `hl.dsp.group.lock()`},
	{GroupLock{Action: ActionEnable}, `hl.dsp.group.lock({ action = "on" })`},

	{GroupLockActive{}, `hl.dsp.group.lock_active()`},
	{GroupLockActive{Action: ActionDisable}, `hl.dsp.group.lock_active({ action = "off" })`},
}

func TestGroupDispatchers(t *testing.T) {
	t.Parallel()

	runCommandTests(t, groupCommandTests)
}

var cursorCommandTests = []commandTest{
	{CursorMoveToCorner{Corner: CornerTopRight}, `hl.dsp.cursor.move_to_corner({ corner = 2 })`},
	{CursorMoveToCorner{Corner: CornerBottomLeft, Window: WindowActive}, `hl.dsp.cursor.move_to_corner({ corner = 0, window = "active" })`},
	{CursorMoveToCorner{Corner: CornerTopLeft}, `hl.dsp.cursor.move_to_corner({ corner = 3 })`},

	{CursorMove{X: 100, Y: 200}, `hl.dsp.cursor.move({ x = 100, y = 200 })`},
	{CursorMove{}, `hl.dsp.cursor.move({ x = 0, y = 0 })`},
}

func TestCursorDispatchers(t *testing.T) {
	t.Parallel()

	runCommandTests(t, cursorCommandTests)
}
