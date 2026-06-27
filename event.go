package hypr

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Event names as Hyprland emits them before ">>". Each has a struct of the same
// name plus the Event suffix.
const (
	EventWorkspace   = "workspace"   // emitted on workspace change. Emitted ONLY when a user requests a workspace change, not on mouse movements (see focusedmon)
	EventWorkspaceV2 = "workspacev2" // emitted on workspace change. Emitted ONLY when a user requests a workspace change, not on mouse movements (see focusedmon)

	EventFocusedMon   = "focusedmon"   // emitted on the active monitor being changed.
	EventFocusedMonV2 = "focusedmonv2" // emitted on the active monitor being changed.

	EventActiveWindow   = "activewindow"   // emitted on the active window being changed.
	EventActiveWindowV2 = "activewindowv2" // emitted on the active window being changed.

	EventFullscreen = "fullscreen" // emitted when a fullscreen status of a window changes.

	EventMonitorRemoved   = "monitorremoved"   // emitted when a monitor is removed (disconnected)
	EventMonitorRemovedV2 = "monitorremovedv2" // emitted when a monitor is removed (disconnected)
	EventMonitorAdded     = "monitoradded"     // emitted when a monitor is added (connected)
	EventMonitorAddedV2   = "monitoraddedv2"   // emitted when a monitor is added (connected)

	EventCreateWorkspace    = "createworkspace"    // emitted when a workspace is created
	EventCreateWorkspaceV2  = "createworkspacev2"  // emitted when a workspace is created
	EventDestroyWorkspace   = "destroyworkspace"   // emitted when a workspace is destroyed
	EventDestroyWorkspaceV2 = "destroyworkspacev2" // emitted when a workspace is destroyed

	EventOpenLayer     = "openlayer"     // emitted when a layerSurface is mapped
	EventCloseLayer    = "closelayer"    // emitted when a layerSurface is unmapped
	EventWindowTitle   = "windowtitle"   // emitted when a window title changes.
	EventWindowTitleV2 = "windowtitlev2" // emitted when a window title changes.
	EventOpenWindow    = "openwindow"    // emitted when a window is opened
	EventCloseWindow   = "closewindow"   // emitted when a window is closed
	EventKill          = "kill"          // emitted when a window is killed (via hyprctl kill)
	EventMoveWindow    = "movewindow"    // emitted when a window is moved to a workspace
	EventMoveWindowV2  = "movewindowv2"  // emitted when a window is moved to a workspace

	EventMoveWorkspace     = "moveworkspace"     // emitted when a workspace is moved to a different monitor
	EventMoveWorkspaceV2   = "moveworkspacev2"   // emitted when a workspace is moved to a different monitor
	EventRenameWorkspace   = "renameworkspace"   // emitted when a workspace is renamed
	EventChangeWorkspaceID = "changeworkspaceid" // emitted when a workspace's ID changes. Returns the old and the new ID.

	EventActiveSpecial   = "activespecial"   // emitted when the special workspace opened in a monitor changes
	EventActiveSpecialV2 = "activespecialv2" // emitted when the special workspace opened in a monitor changes
	EventActiveLayout    = "activelayout"    // emitted on a layout change of the active keyboard

	EventSubmap             = "submap"             // emitted when a keybind submap changes. Empty means default.
	EventChangeFloatingMode = "changefloatingmode" // emitted when a window changes its floating mode. FLOATING is either 0 or 1.
	EventUrgent             = "urgent"             // emitted when a window requests an urgent state
	EventScreenCast         = "screencast"         // emitted when a screencopy state of a client changes. State is 0/1, owner is monitor/window/region.
	EventScreenCastV2       = "screencastv2"       // emitted when a screencopy state of a client changes. State is 0/1, owner is monitor/window/region, name is the identifier of the shared target.

	EventToggleGroup     = "togglegroup"     // emitted when togglegroup command is used. Returns state (0/1) and one or more comma-separated window addresses, e.g. 0,64cea2525760,64cea2522380 where 0 means a group was destroyed and the rest are the windows that were part of it.
	EventMoveIntoGroup   = "moveintogroup"   // emitted when the window is merged into a group. Returns the address of the merged window.
	EventMoveOutOfGroup  = "moveoutofgroup"  // emitted when the window is removed from a group. Returns the address of the removed window.
	EventIgnoreGroupLock = "ignoregrouplock" // emitted when ignoregrouplock is toggled.
	EventLockGroups      = "lockgroups"      // emitted when lockgroups is toggled.

	EventConfigReloaded = "configreloaded" // emitted when the config is done reloading
	EventPin            = "pin"            // emitted when a window is pinned or unpinned
	EventMinimized      = "minimized"      // emitted when an external taskbar-like app requests a window to be minimized
	EventBell           = "bell"           // emitted when an app requests to ring the system bell via xdg-system-bell-v1. Window address parameter may be empty.

	EventCustom = "custom" // emitted by the event dispatcher. Carries whatever string the caller passed, uninterpreted.
)

// ErrInvalidEvent is returned by Parse for a line without ">>".
var ErrInvalidEvent = errors.New("not an event")

var eventSeparator = []byte(">>")

// Event is implemented by every Hyprland event type. Name is the event name as
// Hyprland emits it, "activewindowv2".
type Event interface {
	Name() string
	isEvent()
}

// UnknownEvent is an event this version has no type for, or a known one whose
// payload failed to parse, in which case Err says why. It reaches wildcard
// handlers and the iterator so a newer Hyprland never leaves the caller stuck.
type UnknownEvent struct {
	Line string
	Err  error
}

// Name returns the part of Line before ">>".
func (e *UnknownEvent) Name() string {
	name, _, _ := strings.Cut(e.Line, string(eventSeparator))
	return name
}

// Data returns the part of Line after ">>".
func (e *UnknownEvent) Data() string {
	_, data, _ := strings.Cut(e.Line, string(eventSeparator))
	return data
}

func (*UnknownEvent) isEvent() {}

// Parse turns one event line, with or without its trailing newline, into an
// Event. An unknown name or a payload that does not parse yields an
// *UnknownEvent; only a line without ">>" is an error, ErrInvalidEvent.
func Parse(line []byte) (Event, error) {
	line = bytes.TrimSpace(line)

	name, data, ok := bytes.Cut(line, eventSeparator)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrInvalidEvent, line)
	}

	parse, ok := registry[string(name)]
	if !ok {
		return &UnknownEvent{Line: string(line)}, nil
	}

	event, err := parse(data)
	if err != nil {
		return &UnknownEvent{Line: string(line), Err: err}, nil
	}

	return event, nil
}

// ScreenCastOwner says what a screencast shares.
type ScreenCastOwner int

// Screencast owners as Hyprland emits them.
const (
	ScreenCastMonitor ScreenCastOwner = 0
	ScreenCastWindow  ScreenCastOwner = 1
)

// String returns "monitor" or "window", or "owner(N)" for a value this version
// does not know.
func (o ScreenCastOwner) String() string {
	switch o {
	case ScreenCastMonitor:
		return "monitor"
	case ScreenCastWindow:
		return "window"
	default:
		return "owner(" + strconv.Itoa(int(o)) + ")"
	}
}

// atoi parses an int field. Empty is 0, so an older Hyprland that emits fewer
// fields still parses.
func atoi(field string) (int, error) {
	if field == "" {
		return 0, nil
	}

	return strconv.Atoi(field)
}

// flag parses a bool field, which Hyprland emits as exactly 0 or 1. Empty is
// false, for the same reason atoi treats it as 0.
func flag(field string) (bool, error) {
	switch field {
	case "", "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, fmt.Errorf("flag %q is neither 0 nor 1", field)
	}
}

// payload is one event's data split at commas into fields. The typed accessors
// record the first parse failure in err instead of returning it, so a registry
// entry stays one expression. Missing trailing fields read as zero values.
type payload struct {
	fields []string
	err    error
}

// splitPayload splits data into at most n fields; the last one keeps any
// further commas. n < 0 splits at every comma.
func splitPayload(data []byte, n int) *payload {
	return &payload{fields: strings.SplitN(string(data), ",", n)}
}

func (p *payload) str(i int) string {
	if i < len(p.fields) {
		return p.fields[i]
	}

	return ""
}

func (p *payload) addr(i int) WindowAddress {
	return WindowAddress(p.str(i))
}

func (p *payload) addrs(i int) []WindowAddress {
	if i >= len(p.fields) {
		return nil
	}

	addrs := make([]WindowAddress, 0, len(p.fields)-i)
	for _, f := range p.fields[i:] {
		addrs = append(addrs, WindowAddress(f))
	}

	return addrs
}

func (p *payload) num(i int) int {
	n, err := atoi(p.str(i))
	p.fail(i, err)

	return n
}

func (p *payload) flag(i int) bool {
	b, err := flag(p.str(i))
	p.fail(i, err)

	return b
}

func (p *payload) owner(i int) ScreenCastOwner {
	return ScreenCastOwner(p.num(i))
}

func (p *payload) fail(i int, err error) {
	if err != nil && p.err == nil {
		p.err = fmt.Errorf("field %d: %w", i, err)
	}
}

// registry maps an event name to the parser for its payload, the part after
// ">>". Each parser splits at most as often as the event has fields, so a
// comma inside the last field stays intact.
var registry = map[string]func(data []byte) (Event, error){
	EventActiveWindow: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &ActiveWindowEvent{Class: p.str(0), Title: p.str(1)}, p.err
	},
	EventActiveWindowV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &ActiveWindowV2Event{Address: p.addr(0)}, p.err
	},
	EventWorkspace: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &WorkspaceEvent{WorkspaceName: p.str(0)}, p.err
	},
	EventWorkspaceV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &WorkspaceV2Event{WorkspaceID: p.num(0), WorkspaceName: p.str(1)}, p.err
	},
	EventCreateWorkspace: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &CreateWorkspaceEvent{WorkspaceName: p.str(0)}, p.err
	},
	EventCreateWorkspaceV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &CreateWorkspaceV2Event{WorkspaceID: p.num(0), WorkspaceName: p.str(1)}, p.err
	},
	EventDestroyWorkspace: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &DestroyWorkspaceEvent{WorkspaceName: p.str(0)}, p.err
	},
	EventDestroyWorkspaceV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &DestroyWorkspaceV2Event{WorkspaceID: p.num(0), WorkspaceName: p.str(1)}, p.err
	},
	EventOpenLayer: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &OpenLayerEvent{Namespace: p.str(0)}, p.err
	},
	EventCloseLayer: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &CloseLayerEvent{Namespace: p.str(0)}, p.err
	},
	EventWindowTitle: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &WindowTitleEvent{Address: p.addr(0)}, p.err
	},
	EventWindowTitleV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &WindowTitleV2Event{Address: p.addr(0), Title: p.str(1)}, p.err
	},
	EventOpenWindow: func(data []byte) (Event, error) {
		p := splitPayload(data, 4)
		return &OpenWindowEvent{Address: p.addr(0), WorkspaceName: p.str(1), Class: p.str(2), Title: p.str(3)}, p.err
	},
	EventCloseWindow: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &CloseWindowEvent{Address: p.addr(0)}, p.err
	},
	EventFocusedMon: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &FocusedMonEvent{MonitorName: p.str(0), WorkspaceName: p.str(1)}, p.err
	},
	EventFocusedMonV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &FocusedMonV2Event{MonitorName: p.str(0), WorkspaceID: p.num(1)}, p.err
	},
	EventFullscreen: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &FullscreenEvent{Fullscreen: p.flag(0)}, p.err
	},
	EventMonitorRemoved: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &MonitorRemovedEvent{MonitorName: p.str(0)}, p.err
	},
	EventMonitorRemovedV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 3)
		return &MonitorRemovedV2Event{MonitorID: p.num(0), MonitorName: p.str(1), MonitorDescription: p.str(2)}, p.err
	},
	EventMonitorAdded: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &MonitorAddedEvent{MonitorName: p.str(0)}, p.err
	},
	EventMonitorAddedV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 3)
		return &MonitorAddedV2Event{MonitorID: p.num(0), MonitorName: p.str(1), MonitorDescription: p.str(2)}, p.err
	},
	EventKill: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &KillEvent{Address: p.addr(0)}, p.err
	},
	EventMoveWindow: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &MoveWindowEvent{Address: p.addr(0), WorkspaceName: p.str(1)}, p.err
	},
	EventMoveWindowV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 3)
		return &MoveWindowV2Event{Address: p.addr(0), WorkspaceID: p.num(1), WorkspaceName: p.str(2)}, p.err
	},
	EventMoveWorkspace: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &MoveWorkspaceEvent{WorkspaceName: p.str(0), MonitorName: p.str(1)}, p.err
	},
	EventMoveWorkspaceV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 3)
		return &MoveWorkspaceV2Event{WorkspaceID: p.num(0), WorkspaceName: p.str(1), MonitorName: p.str(2)}, p.err
	},
	EventRenameWorkspace: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &RenameWorkspaceEvent{WorkspaceID: p.num(0), NewName: p.str(1)}, p.err
	},
	EventChangeWorkspaceID: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &ChangeWorkspaceIDEvent{OldID: p.num(0), NewID: p.num(1)}, p.err
	},
	EventActiveSpecial: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &ActiveSpecialEvent{WorkspaceName: p.str(0), MonitorName: p.str(1)}, p.err
	},
	EventActiveSpecialV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 3)
		return &ActiveSpecialV2Event{WorkspaceID: p.num(0), WorkspaceName: p.str(1), MonitorName: p.str(2)}, p.err
	},
	EventActiveLayout: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &ActiveLayoutEvent{KeyboardName: p.str(0), LayoutName: p.str(1)}, p.err
	},
	EventSubmap: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &SubmapEvent{Submap: p.str(0)}, p.err
	},
	EventChangeFloatingMode: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &ChangeFloatingModeEvent{Address: p.addr(0), Floating: p.flag(1)}, p.err
	},
	EventUrgent: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &UrgentEvent{Address: p.addr(0)}, p.err
	},
	EventScreenCast: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &ScreenCastEvent{Sharing: p.flag(0), Owner: p.owner(1)}, p.err
	},
	EventScreenCastV2: func(data []byte) (Event, error) {
		p := splitPayload(data, 3)
		return &ScreenCastV2Event{Sharing: p.flag(0), Owner: p.owner(1), TargetName: p.str(2)}, p.err
	},
	EventToggleGroup: func(data []byte) (Event, error) {
		p := splitPayload(data, -1)
		return &ToggleGroupEvent{Grouped: p.flag(0), Addresses: p.addrs(1)}, p.err
	},
	EventMoveIntoGroup: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &MoveIntoGroupEvent{Address: p.addr(0)}, p.err
	},
	EventMoveOutOfGroup: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &MoveOutOfGroupEvent{Address: p.addr(0)}, p.err
	},
	EventIgnoreGroupLock: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &IgnoreGroupLockEvent{Ignored: p.flag(0)}, p.err
	},
	EventLockGroups: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &LockGroupsEvent{Locked: p.flag(0)}, p.err
	},
	EventConfigReloaded: func(_ []byte) (Event, error) {
		return &ConfigReloadedEvent{}, nil
	},
	EventPin: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &PinEvent{Address: p.addr(0), Pinned: p.flag(1)}, p.err
	},
	EventMinimized: func(data []byte) (Event, error) {
		p := splitPayload(data, 2)
		return &MinimizedEvent{Address: p.addr(0), Minimized: p.flag(1)}, p.err
	},
	EventBell: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &BellEvent{Address: p.addr(0)}, p.err
	},
	EventCustom: func(data []byte) (Event, error) {
		p := splitPayload(data, 1)
		return &CustomEvent{Data: p.str(0)}, p.err
	},
}
