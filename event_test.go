package hypr

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

// Every event type must satisfy Event through its pointer. Both methods have
// pointer receivers, so the value type does not.
var (
	_ Event = (*ActiveWindowEvent)(nil)
	_ Event = (*ActiveWindowV2Event)(nil)
	_ Event = (*WorkspaceEvent)(nil)
	_ Event = (*WorkspaceV2Event)(nil)
	_ Event = (*CreateWorkspaceEvent)(nil)
	_ Event = (*CreateWorkspaceV2Event)(nil)
	_ Event = (*DestroyWorkspaceEvent)(nil)
	_ Event = (*DestroyWorkspaceV2Event)(nil)
	_ Event = (*OpenLayerEvent)(nil)
	_ Event = (*CloseLayerEvent)(nil)
	_ Event = (*WindowTitleEvent)(nil)
	_ Event = (*WindowTitleV2Event)(nil)
	_ Event = (*OpenWindowEvent)(nil)
	_ Event = (*CloseWindowEvent)(nil)
	_ Event = (*FocusedMonEvent)(nil)
	_ Event = (*FocusedMonV2Event)(nil)
	_ Event = (*FullscreenEvent)(nil)
	_ Event = (*MonitorRemovedEvent)(nil)
	_ Event = (*MonitorRemovedV2Event)(nil)
	_ Event = (*MonitorAddedEvent)(nil)
	_ Event = (*MonitorAddedV2Event)(nil)
	_ Event = (*KillEvent)(nil)
	_ Event = (*MoveWindowEvent)(nil)
	_ Event = (*MoveWindowV2Event)(nil)
	_ Event = (*MoveWorkspaceEvent)(nil)
	_ Event = (*MoveWorkspaceV2Event)(nil)
	_ Event = (*RenameWorkspaceEvent)(nil)
	_ Event = (*ActiveSpecialEvent)(nil)
	_ Event = (*ActiveSpecialV2Event)(nil)
	_ Event = (*ActiveLayoutEvent)(nil)
	_ Event = (*SubmapEvent)(nil)
	_ Event = (*ChangeFloatingModeEvent)(nil)
	_ Event = (*UrgentEvent)(nil)
	_ Event = (*ScreenCastEvent)(nil)
	_ Event = (*ScreenCastV2Event)(nil)
	_ Event = (*ToggleGroupEvent)(nil)
	_ Event = (*MoveIntoGroupEvent)(nil)
	_ Event = (*MoveOutOfGroupEvent)(nil)
	_ Event = (*IgnoreGroupLockEvent)(nil)
	_ Event = (*LockGroupsEvent)(nil)
	_ Event = (*ConfigReloadedEvent)(nil)
	_ Event = (*PinEvent)(nil)
	_ Event = (*MinimizedEvent)(nil)
	_ Event = (*BellEvent)(nil)
	_ Event = (*ChangeWorkspaceIDEvent)(nil)
	_ Event = (*CustomEvent)(nil)
	_ Event = (*UnknownEvent)(nil)
)

// allEventNames lists every event name constant. TestEventNameCoverage keeps it
// in step with the registry and with testdata/events.txt.
var allEventNames = []string{
	EventWorkspace, EventWorkspaceV2,
	EventFocusedMon, EventFocusedMonV2,
	EventActiveWindow, EventActiveWindowV2,
	EventFullscreen,
	EventMonitorRemoved, EventMonitorRemovedV2,
	EventMonitorAdded, EventMonitorAddedV2,
	EventCreateWorkspace, EventCreateWorkspaceV2,
	EventDestroyWorkspace, EventDestroyWorkspaceV2,
	EventOpenLayer, EventCloseLayer,
	EventWindowTitle, EventWindowTitleV2,
	EventOpenWindow, EventCloseWindow,
	EventKill,
	EventMoveWindow, EventMoveWindowV2,
	EventMoveWorkspace, EventMoveWorkspaceV2,
	EventRenameWorkspace, EventChangeWorkspaceID,
	EventActiveSpecial, EventActiveSpecialV2,
	EventActiveLayout,
	EventSubmap, EventChangeFloatingMode, EventUrgent,
	EventScreenCast, EventScreenCastV2,
	EventToggleGroup, EventMoveIntoGroup, EventMoveOutOfGroup,
	EventIgnoreGroupLock, EventLockGroups,
	EventConfigReloaded, EventPin, EventMinimized, EventBell,
	EventCustom,
}

func TestParseNameAndData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       string
		wantName string
		wantData string
		wantErr  error
	}{
		{name: "plain", in: "workspace>>2", wantName: "workspace", wantData: "2"},
		{name: "trailing newline", in: "workspace>>2\n", wantName: "workspace", wantData: "2"},
		{name: "surrounding whitespace", in: "  workspace>>2  \n", wantName: "workspace", wantData: "2"},
		{name: "empty payload", in: "workspace>>", wantName: "workspace"},
		{name: "payload contains separator", in: "submap>>a>>b", wantName: "submap", wantData: "a>>b"},
		{name: "unknown name", in: "notanevent>>2", wantName: "notanevent", wantData: "2"},
		{name: "separator at start", in: ">>2", wantName: "", wantData: "2"},
		{name: "no separator", in: "workspace", wantErr: ErrInvalidEvent},
		{name: "single angle bracket", in: "workspace>2", wantErr: ErrInvalidEvent},
		{name: "empty", in: "", wantErr: ErrInvalidEvent},
		{name: "whitespace only", in: "   \n", wantErr: ErrInvalidEvent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse([]byte(tt.in))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Parse(%q) error = %v, want %v", tt.in, err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if got != nil {
					t.Errorf("Parse(%q) = %+v, want nil", tt.in, got)
				}
				return
			}
			if got.Name() != tt.wantName {
				t.Errorf("Parse(%q).Name() = %q, want %q", tt.in, got.Name(), tt.wantName)
			}

			// The data half is only observable through UnknownEvent.
			if u, ok := got.(*UnknownEvent); ok {
				if u.Data() != tt.wantData {
					t.Errorf("Parse(%q).Data() = %q, want %q", tt.in, u.Data(), tt.wantData)
				}
				if u.Err != nil {
					t.Errorf("Parse(%q).Err = %v, want nil", tt.in, u.Err)
				}
			}
		})
	}
}

// TestParseEvent covers every event name three ways: all fields present,
// trailing fields missing, and commas inside the final field. The SplitN limits
// in the registry exist to keep the last field intact, so that third shape is
// where a wrong limit shows up.
func TestParseEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line string
		want Event
	}{
		// activewindow: class,title
		{"activewindow>>firefox,Home", &ActiveWindowEvent{Class: "firefox", Title: "Home"}},
		{"activewindow>>,", &ActiveWindowEvent{}},
		{"activewindow>>firefox,a,b,c", &ActiveWindowEvent{Class: "firefox", Title: "a,b,c"}},

		// activewindowv2: address
		{"activewindowv2>>56157f7c25a0", &ActiveWindowV2Event{Address: "56157f7c25a0"}},
		{"activewindowv2>>", &ActiveWindowV2Event{}},
		{"activewindowv2>>a,b", &ActiveWindowV2Event{Address: "a,b"}},

		// workspace: name
		{"workspace>>2", &WorkspaceEvent{WorkspaceName: "2"}},
		{"workspace>>", &WorkspaceEvent{}},
		{"workspace>>special:magic", &WorkspaceEvent{WorkspaceName: "special:magic"}},

		// workspacev2: id,name
		{"workspacev2>>2,two", &WorkspaceV2Event{WorkspaceID: 2, WorkspaceName: "two"}},
		{"workspacev2>>2", &WorkspaceV2Event{WorkspaceID: 2}},
		{"workspacev2>>-98,a,b", &WorkspaceV2Event{WorkspaceID: -98, WorkspaceName: "a,b"}},

		// createworkspace: name
		{"createworkspace>>3", &CreateWorkspaceEvent{WorkspaceName: "3"}},
		{"createworkspace>>", &CreateWorkspaceEvent{}},
		{"createworkspace>>a,b", &CreateWorkspaceEvent{WorkspaceName: "a,b"}},

		// createworkspacev2: id,name
		{"createworkspacev2>>3,three", &CreateWorkspaceV2Event{WorkspaceID: 3, WorkspaceName: "three"}},
		{"createworkspacev2>>3", &CreateWorkspaceV2Event{WorkspaceID: 3}},
		{"createworkspacev2>>3,a,b", &CreateWorkspaceV2Event{WorkspaceID: 3, WorkspaceName: "a,b"}},

		// destroyworkspace: name
		{"destroyworkspace>>3", &DestroyWorkspaceEvent{WorkspaceName: "3"}},
		{"destroyworkspace>>", &DestroyWorkspaceEvent{}},
		{"destroyworkspace>>a,b", &DestroyWorkspaceEvent{WorkspaceName: "a,b"}},

		// destroyworkspacev2: id,name
		{"destroyworkspacev2>>3,three", &DestroyWorkspaceV2Event{WorkspaceID: 3, WorkspaceName: "three"}},
		{"destroyworkspacev2>>3", &DestroyWorkspaceV2Event{WorkspaceID: 3}},
		{"destroyworkspacev2>>3,a,b", &DestroyWorkspaceV2Event{WorkspaceID: 3, WorkspaceName: "a,b"}},

		// openlayer: namespace
		{"openlayer>>waybar", &OpenLayerEvent{Namespace: "waybar"}},
		{"openlayer>>", &OpenLayerEvent{}},
		{"openlayer>>a,b", &OpenLayerEvent{Namespace: "a,b"}},

		// closelayer: namespace
		{"closelayer>>waybar", &CloseLayerEvent{Namespace: "waybar"}},
		{"closelayer>>", &CloseLayerEvent{}},
		{"closelayer>>a,b", &CloseLayerEvent{Namespace: "a,b"}},

		// windowtitle: address
		{"windowtitle>>56157f7c25a0", &WindowTitleEvent{Address: "56157f7c25a0"}},
		{"windowtitle>>", &WindowTitleEvent{}},
		{"windowtitle>>a,b", &WindowTitleEvent{Address: "a,b"}},

		// windowtitlev2: address,title
		{"windowtitlev2>>56157f7c25a0,zsh", &WindowTitleV2Event{Address: "56157f7c25a0", Title: "zsh"}},
		{"windowtitlev2>>56157f7c25a0", &WindowTitleV2Event{Address: "56157f7c25a0"}},
		{"windowtitlev2>>56157f7c25a0,a,b", &WindowTitleV2Event{Address: "56157f7c25a0", Title: "a,b"}},

		// openwindow: address,workspace,class,title
		{"openwindow>>56157e75fd30,2,firefox,Mozilla Firefox", &OpenWindowEvent{Address: "56157e75fd30", WorkspaceName: "2", Class: "firefox", Title: "Mozilla Firefox"}},
		{"openwindow>>56157e75fd30,2", &OpenWindowEvent{Address: "56157e75fd30", WorkspaceName: "2"}},
		{"openwindow>>56157e75fd30,2,firefox,a,b", &OpenWindowEvent{Address: "56157e75fd30", WorkspaceName: "2", Class: "firefox", Title: "a,b"}},

		// closewindow: address
		{"closewindow>>56157f7c25a0", &CloseWindowEvent{Address: "56157f7c25a0"}},
		{"closewindow>>", &CloseWindowEvent{}},
		{"closewindow>>a,b", &CloseWindowEvent{Address: "a,b"}},

		// focusedmon: monname,workspacename
		{"focusedmon>>DP-1,2", &FocusedMonEvent{MonitorName: "DP-1", WorkspaceName: "2"}},
		{"focusedmon>>DP-1", &FocusedMonEvent{MonitorName: "DP-1"}},
		{"focusedmon>>DP-1,a,b", &FocusedMonEvent{MonitorName: "DP-1", WorkspaceName: "a,b"}},

		// focusedmonv2: monname,workspaceid
		{"focusedmonv2>>DP-1,2", &FocusedMonV2Event{MonitorName: "DP-1", WorkspaceID: 2}},
		{"focusedmonv2>>DP-1", &FocusedMonV2Event{MonitorName: "DP-1"}},
		{"focusedmonv2>>DP-1,", &FocusedMonV2Event{MonitorName: "DP-1"}},

		// fullscreen: 0/1
		{"fullscreen>>1", &FullscreenEvent{Fullscreen: true}},
		{"fullscreen>>0", &FullscreenEvent{}},
		{"fullscreen>>", &FullscreenEvent{}},

		// monitorremoved: name
		{"monitorremoved>>DP-1", &MonitorRemovedEvent{MonitorName: "DP-1"}},
		{"monitorremoved>>", &MonitorRemovedEvent{}},
		{"monitorremoved>>a,b", &MonitorRemovedEvent{MonitorName: "a,b"}},

		// monitorremovedv2: id,name,description
		{"monitorremovedv2>>0,DP-1,Dell U2720Q", &MonitorRemovedV2Event{MonitorID: 0, MonitorName: "DP-1", MonitorDescription: "Dell U2720Q"}},
		{"monitorremovedv2>>1,DP-1", &MonitorRemovedV2Event{MonitorID: 1, MonitorName: "DP-1"}},
		{"monitorremovedv2>>1,DP-1,Dell, Inc.", &MonitorRemovedV2Event{MonitorID: 1, MonitorName: "DP-1", MonitorDescription: "Dell, Inc."}},

		// monitoradded: name
		{"monitoradded>>DP-1", &MonitorAddedEvent{MonitorName: "DP-1"}},
		{"monitoradded>>", &MonitorAddedEvent{}},
		{"monitoradded>>a,b", &MonitorAddedEvent{MonitorName: "a,b"}},

		// monitoraddedv2: id,name,description
		{"monitoraddedv2>>0,DP-1,Dell U2720Q", &MonitorAddedV2Event{MonitorID: 0, MonitorName: "DP-1", MonitorDescription: "Dell U2720Q"}},
		{"monitoraddedv2>>1,DP-1", &MonitorAddedV2Event{MonitorID: 1, MonitorName: "DP-1"}},
		{"monitoraddedv2>>1,DP-1,Dell, Inc.", &MonitorAddedV2Event{MonitorID: 1, MonitorName: "DP-1", MonitorDescription: "Dell, Inc."}},

		// kill: address
		{"kill>>56157f7c25a0", &KillEvent{Address: "56157f7c25a0"}},
		{"kill>>", &KillEvent{}},
		{"kill>>a,b", &KillEvent{Address: "a,b"}},

		// movewindow: address,workspacename
		{"movewindow>>56157f7c25a0,2", &MoveWindowEvent{Address: "56157f7c25a0", WorkspaceName: "2"}},
		{"movewindow>>56157f7c25a0", &MoveWindowEvent{Address: "56157f7c25a0"}},
		{"movewindow>>56157f7c25a0,a,b", &MoveWindowEvent{Address: "56157f7c25a0", WorkspaceName: "a,b"}},

		// movewindowv2: address,workspaceid,workspacename
		{"movewindowv2>>56157f7c25a0,2,two", &MoveWindowV2Event{Address: "56157f7c25a0", WorkspaceID: 2, WorkspaceName: "two"}},
		{"movewindowv2>>56157f7c25a0,2", &MoveWindowV2Event{Address: "56157f7c25a0", WorkspaceID: 2}},
		{"movewindowv2>>56157f7c25a0,2,a,b", &MoveWindowV2Event{Address: "56157f7c25a0", WorkspaceID: 2, WorkspaceName: "a,b"}},

		// moveworkspace: workspacename,monname
		{"moveworkspace>>2,DP-1", &MoveWorkspaceEvent{WorkspaceName: "2", MonitorName: "DP-1"}},
		{"moveworkspace>>2", &MoveWorkspaceEvent{WorkspaceName: "2"}},
		{"moveworkspace>>2,a,b", &MoveWorkspaceEvent{WorkspaceName: "2", MonitorName: "a,b"}},

		// moveworkspacev2: workspaceid,workspacename,monname
		{"moveworkspacev2>>2,two,DP-1", &MoveWorkspaceV2Event{WorkspaceID: 2, WorkspaceName: "two", MonitorName: "DP-1"}},
		{"moveworkspacev2>>2,two", &MoveWorkspaceV2Event{WorkspaceID: 2, WorkspaceName: "two"}},
		{"moveworkspacev2>>2,two,a,b", &MoveWorkspaceV2Event{WorkspaceID: 2, WorkspaceName: "two", MonitorName: "a,b"}},

		// renameworkspace: workspaceid,newname
		{"renameworkspace>>2,web", &RenameWorkspaceEvent{WorkspaceID: 2, NewName: "web"}},
		{"renameworkspace>>2", &RenameWorkspaceEvent{WorkspaceID: 2}},
		{"renameworkspace>>2,a,b", &RenameWorkspaceEvent{WorkspaceID: 2, NewName: "a,b"}},

		// changeworkspaceid: oldid,newid
		{"changeworkspaceid>>6,77", &ChangeWorkspaceIDEvent{OldID: 6, NewID: 77}},
		{"changeworkspaceid>>6", &ChangeWorkspaceIDEvent{OldID: 6}},
		{"changeworkspaceid>>6,", &ChangeWorkspaceIDEvent{OldID: 6}},

		// activespecial: workspacename,monname
		{"activespecial>>special:magic,DP-1", &ActiveSpecialEvent{WorkspaceName: "special:magic", MonitorName: "DP-1"}},
		{"activespecial>>,DP-1", &ActiveSpecialEvent{MonitorName: "DP-1"}},
		{"activespecial>>special:magic,a,b", &ActiveSpecialEvent{WorkspaceName: "special:magic", MonitorName: "a,b"}},

		// activespecialv2: workspaceid,workspacename,monname
		{"activespecialv2>>-98,special:magic,DP-1", &ActiveSpecialV2Event{WorkspaceID: -98, WorkspaceName: "special:magic", MonitorName: "DP-1"}},
		{"activespecialv2>>,,DP-1", &ActiveSpecialV2Event{MonitorName: "DP-1"}},
		{"activespecialv2>>-98,special:magic,a,b", &ActiveSpecialV2Event{WorkspaceID: -98, WorkspaceName: "special:magic", MonitorName: "a,b"}},

		// activelayout: keyboardname,layoutname
		{"activelayout>>at-translated-set-2-keyboard,English (US)", &ActiveLayoutEvent{KeyboardName: "at-translated-set-2-keyboard", LayoutName: "English (US)"}},
		{"activelayout>>kbd", &ActiveLayoutEvent{KeyboardName: "kbd"}},
		{"activelayout>>kbd,a,b", &ActiveLayoutEvent{KeyboardName: "kbd", LayoutName: "a,b"}},

		// submap: name
		{"submap>>resize", &SubmapEvent{Submap: "resize"}},
		{"submap>>", &SubmapEvent{}},
		{"submap>>a,b", &SubmapEvent{Submap: "a,b"}},

		// changefloatingmode: address,floating
		{"changefloatingmode>>56157f7c25a0,1", &ChangeFloatingModeEvent{Address: "56157f7c25a0", Floating: true}},
		{"changefloatingmode>>56157f7c25a0", &ChangeFloatingModeEvent{Address: "56157f7c25a0"}},
		{"changefloatingmode>>56157f7c25a0,0", &ChangeFloatingModeEvent{Address: "56157f7c25a0"}},

		// urgent: address
		{"urgent>>56157f7c25a0", &UrgentEvent{Address: "56157f7c25a0"}},
		{"urgent>>", &UrgentEvent{}},
		{"urgent>>a,b", &UrgentEvent{Address: "a,b"}},

		// screencast: state,owner
		{"screencast>>1,0", &ScreenCastEvent{Sharing: true, Owner: ScreenCastMonitor}},
		{"screencast>>1", &ScreenCastEvent{Sharing: true}},
		{"screencast>>0,1", &ScreenCastEvent{Owner: ScreenCastWindow}},

		// screencastv2: state,owner,name
		{"screencastv2>>1,1,56157f7c25a0", &ScreenCastV2Event{Sharing: true, Owner: ScreenCastWindow, TargetName: "56157f7c25a0"}},
		{"screencastv2>>1,0", &ScreenCastV2Event{Sharing: true, Owner: ScreenCastMonitor}},
		{"screencastv2>>0,0,a,b", &ScreenCastV2Event{TargetName: "a,b"}},

		// togglegroup: state,address[,address...]
		{"togglegroup>>1,64cea2525760", &ToggleGroupEvent{Grouped: true, Addresses: []WindowAddress{"64cea2525760"}}},
		{"togglegroup>>0", &ToggleGroupEvent{}},
		{"togglegroup>>0,64cea2525760,64cea2522380", &ToggleGroupEvent{Addresses: []WindowAddress{"64cea2525760", "64cea2522380"}}},

		// moveintogroup: address
		{"moveintogroup>>56157f7c25a0", &MoveIntoGroupEvent{Address: "56157f7c25a0"}},
		{"moveintogroup>>", &MoveIntoGroupEvent{}},
		{"moveintogroup>>a,b", &MoveIntoGroupEvent{Address: "a,b"}},

		// moveoutofgroup: address
		{"moveoutofgroup>>56157f7c25a0", &MoveOutOfGroupEvent{Address: "56157f7c25a0"}},
		{"moveoutofgroup>>", &MoveOutOfGroupEvent{}},
		{"moveoutofgroup>>a,b", &MoveOutOfGroupEvent{Address: "a,b"}},

		// ignoregrouplock: 0/1
		{"ignoregrouplock>>1", &IgnoreGroupLockEvent{Ignored: true}},
		{"ignoregrouplock>>0", &IgnoreGroupLockEvent{}},
		{"ignoregrouplock>>", &IgnoreGroupLockEvent{}},

		// lockgroups: 0/1
		{"lockgroups>>1", &LockGroupsEvent{Locked: true}},
		{"lockgroups>>0", &LockGroupsEvent{}},
		{"lockgroups>>", &LockGroupsEvent{}},

		// configreloaded: no params
		{"configreloaded>>", &ConfigReloadedEvent{}},
		{"configreloaded>>x", &ConfigReloadedEvent{}},
		{"configreloaded>>a,b", &ConfigReloadedEvent{}},

		// pin: address,pinstate
		{"pin>>56157f7c25a0,1", &PinEvent{Address: "56157f7c25a0", Pinned: true}},
		{"pin>>56157f7c25a0", &PinEvent{Address: "56157f7c25a0"}},
		{"pin>>56157f7c25a0,0", &PinEvent{Address: "56157f7c25a0"}},

		// minimized: address,state
		{"minimized>>56157f7c25a0,1", &MinimizedEvent{Address: "56157f7c25a0", Minimized: true}},
		{"minimized>>56157f7c25a0", &MinimizedEvent{Address: "56157f7c25a0"}},
		{"minimized>>56157f7c25a0,0", &MinimizedEvent{Address: "56157f7c25a0"}},

		// bell: address (may be empty)
		{"bell>>56157f7c25a0", &BellEvent{Address: "56157f7c25a0"}},
		{"bell>>", &BellEvent{}},
		{"bell>>a,b", &BellEvent{Address: "a,b"}},

		// custom: an arbitrary string, not split
		{"custom>>hello, world", &CustomEvent{Data: "hello, world"}},
		{"custom>>", &CustomEvent{}},
		{"custom>>a>>b", &CustomEvent{Data: "a>>b"}},
	}

	seen := make(map[string]bool)
	for _, tt := range tests {
		name, _, _ := strings.Cut(tt.line, ">>")
		seen[name] = true
		checkEventTypeName(t, name, tt.want)
		t.Run(tt.line, func(t *testing.T) {
			t.Parallel()

			got, err := Parse([]byte(tt.line))
			if err != nil {
				t.Fatalf("Parse(%q) = %v", tt.line, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%q)\n got %+v\nwant %+v", tt.line, got, tt.want)
			}
			if name, _, _ := strings.Cut(tt.line, ">>"); got.Name() != name {
				t.Errorf("Name() = %q, want %q", got.Name(), name)
			}
		})
	}

	// Every registered event has a row above, and its Go type is the event
	// name plus Event, see doc.go.
	for name := range registry {
		if !seen[name] {
			t.Errorf("event %q has no row in TestParseEvent", name)
		}
	}
}

// checkEventTypeName enforces the rule in doc.go: the Go type of an event is
// its name plus Event, case aside.
func checkEventTypeName(t *testing.T, name string, event Event) {
	t.Helper()

	typeName := reflect.TypeOf(event).Elem().Name()
	if !strings.HasSuffix(typeName, "Event") || !strings.EqualFold(strings.TrimSuffix(typeName, "Event"), name) {
		t.Errorf("event %q is %s, want a type named %sEvent", name, typeName, name)
	}
}

// TestParseEventTrimsWhitespace pins that a line still parses once bufio hands
// it over with its trailing newline attached.
func TestParseEventTrimsWhitespace(t *testing.T) {
	t.Parallel()

	got, err := Parse([]byte("  workspacev2>>2,two  \n"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	want := &WorkspaceV2Event{WorkspaceID: 2, WorkspaceName: "two"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

// TestParseTypedFields pins the parse rules for int and bool fields: empty is
// the zero value, anything else must parse, and a bool is exactly 0 or 1.
func TestParseTypedFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line string
		want Event
		bad  bool
	}{
		{line: "workspacev2>>,web", want: &WorkspaceV2Event{WorkspaceName: "web"}},
		{line: "workspacev2>>x,web", bad: true},
		{line: "workspacev2>>2.5,web", bad: true},
		{line: "changeworkspaceid>>6,x", bad: true},
		{line: "fullscreen>>2", bad: true},
		{line: "fullscreen>>true", bad: true},
		{line: "fullscreen>>1", want: &FullscreenEvent{Fullscreen: true}},
		{line: "pin>>addr,x", bad: true},
		{line: "screencast>>1,7", want: &ScreenCastEvent{Sharing: true, Owner: ScreenCastOwner(7)}},
		{line: "screencast>>1,x", bad: true},
		{line: "togglegroup>>2,addr", bad: true},
		{line: "togglegroup>>1,", want: &ToggleGroupEvent{Grouped: true, Addresses: []WindowAddress{""}}},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			t.Parallel()

			got, err := Parse([]byte(tt.line))
			if err != nil {
				t.Fatalf("Parse(%q) = %v", tt.line, err)
			}

			u, unknown := got.(*UnknownEvent)
			if tt.bad {
				if !unknown {
					t.Fatalf("Parse(%q) = %+v, want *UnknownEvent", tt.line, got)
				}
				if u.Err == nil {
					t.Errorf("Parse(%q).Err = nil, want the parse failure", tt.line)
				}
				if u.Line != tt.line {
					t.Errorf("Parse(%q).Line = %q", tt.line, u.Line)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%q)\n got %+v\nwant %+v", tt.line, got, tt.want)
			}
		})
	}
}

func TestScreenCastOwnerString(t *testing.T) {
	t.Parallel()

	tests := map[ScreenCastOwner]string{
		ScreenCastMonitor:  "monitor",
		ScreenCastWindow:   "window",
		ScreenCastOwner(7): "owner(7)",
	}
	for owner, want := range tests {
		if got := owner.String(); got != want {
			t.Errorf("ScreenCastOwner(%d).String() = %q, want %q", int(owner), got, want)
		}
	}
}

func TestParseUnknown(t *testing.T) {
	t.Parallel()

	got, err := Parse([]byte("notanevent>>2,3\n"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	want := &UnknownEvent{Line: "notanevent>>2,3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse = %+v, want %+v", got, want)
	}
	if got.Name() != "notanevent" {
		t.Errorf("Name() = %q, want notanevent", got.Name())
	}
	if data := got.(*UnknownEvent).Data(); data != "2,3" {
		t.Errorf("Data() = %q, want 2,3", data)
	}
}

// TestParseEventStream replays a captured Hyprland socket dump. It carries the
// awkward inputs a synthetic table tends to miss: blank lines between bursts and
// empty payloads such as "activewindow>>," when focus lands on nothing.
func TestParseEventStream(t *testing.T) {
	t.Parallel()

	var parsed int
	for i, line := range readLines(t, "testdata/event_stream.txt") {
		if strings.TrimSpace(line) == "" {
			// Blank lines carry no event and must not reach a handler.
			if _, err := Parse([]byte(line)); !errors.Is(err, ErrInvalidEvent) {
				t.Errorf("line %d: blank line error = %v, want ErrInvalidEvent", i+1, err)
			}
			continue
		}

		event, err := Parse([]byte(line))
		if err != nil {
			t.Errorf("line %d: Parse(%q) = %v", i+1, line, err)
			continue
		}
		if u, ok := event.(*UnknownEvent); ok {
			t.Errorf("line %d: Parse(%q) = UnknownEvent, err %v", i+1, line, u.Err)
			continue
		}
		if name, _, _ := strings.Cut(line, ">>"); event.Name() != name {
			t.Errorf("line %d: Name() = %q, want %q", i+1, event.Name(), name)
		}
		parsed++
	}

	if parsed == 0 {
		t.Fatal("no events parsed from the stream corpus")
	}
	t.Logf("parsed %d events", parsed)
}

// TestEventNameCoverage keeps the constants in step with the event names
// Hyprland documents, listed in testdata/events.txt.
func TestEventNameCoverage(t *testing.T) {
	t.Parallel()

	// Documented by Hyprland but not modelled yet. Remove an entry here once it
	// gets a constant and a registry entry.
	knownMissing := map[string]bool{}
	// Handled by us but absent from testdata/events.txt.
	knownExtra := map[string]bool{
		EventIgnoreGroupLock: true,
	}

	handled := make(map[string]bool, len(allEventNames))
	for _, name := range allEventNames {
		if handled[name] {
			t.Errorf("event %q listed twice in allEventNames", name)
		}
		handled[name] = true

		if _, ok := registry[name]; !ok {
			t.Errorf("event %q has a constant but no registry entry", name)
			continue
		}
		event, err := Parse([]byte(name + ">>x"))
		if err != nil {
			t.Errorf("Parse(%q) = %v", name, err)
			continue
		}
		if event.Name() != name {
			t.Errorf("Parse(%q).Name() = %q", name, event.Name())
		}
	}
	for name := range registry {
		if !handled[name] {
			t.Errorf("event %q is in the registry but not in allEventNames", name)
		}
	}

	documented := make(map[string]bool)
	for _, name := range readLines(t, "testdata/events.txt") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		documented[name] = true

		if !handled[name] && !knownMissing[name] {
			t.Errorf("event %q is documented but not handled", name)
		}
	}

	for name := range handled {
		if !documented[name] && !knownExtra[name] {
			t.Errorf("event %q is handled but not documented in testdata/events.txt", name)
		}
	}
	for name := range knownMissing {
		if handled[name] {
			t.Errorf("event %q is handled now, drop it from knownMissing", name)
		}
	}
}

func FuzzParseEvent(f *testing.F) {
	for _, line := range readLines(f, "testdata/event_stream.txt") {
		f.Add(line)
	}
	f.Add("workspace")
	f.Add(">>")
	f.Add("togglegroup>>")

	f.Fuzz(func(t *testing.T, line string) {
		event, err := Parse([]byte(line))
		if err != nil {
			return
		}
		if name, _, _ := strings.Cut(strings.TrimSpace(line), ">>"); event.Name() != name {
			t.Errorf("Name() = %q, want %q", event.Name(), name)
		}
	})
}

func readLines(tb testing.TB, path string) []string {
	tb.Helper()

	f, err := os.Open(path)
	if err != nil {
		tb.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		tb.Fatalf("read %s: %v", path, err)
	}
	return lines
}
