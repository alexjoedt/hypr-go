package hypr

// ActiveWindowEvent is emitted when the active window changes. Class and Title
// are both empty when nothing is focused.
//
// Wire: activewindow>>firefox,Mozilla Firefox
type ActiveWindowEvent struct {
	Class string
	Title string
}

// Name returns EventActiveWindow.
func (*ActiveWindowEvent) Name() string { return EventActiveWindow }

func (*ActiveWindowEvent) isEvent() {}

// ActiveWindowV2Event is the v2 variant, providing only the window address.
// Address is empty when nothing is focused.
//
// Wire: activewindowv2>>56157f7c25a0
type ActiveWindowV2Event struct {
	Address WindowAddress
}

// Name returns EventActiveWindowV2.
func (*ActiveWindowV2Event) Name() string { return EventActiveWindowV2 }

func (*ActiveWindowV2Event) isEvent() {}

// WorkspaceEvent is emitted when the active workspace changes. For numbered
// workspaces the name is the number.
//
// Wire: workspace>>2
type WorkspaceEvent struct {
	WorkspaceName string
}

// Name returns EventWorkspace.
func (*WorkspaceEvent) Name() string { return EventWorkspace }

func (*WorkspaceEvent) isEvent() {}

// WorkspaceV2Event is the v2 variant, providing the workspace ID as well.
//
// Wire: workspacev2>>2,web
type WorkspaceV2Event struct {
	WorkspaceID   int
	WorkspaceName string
}

// Name returns EventWorkspaceV2.
func (*WorkspaceV2Event) Name() string { return EventWorkspaceV2 }

func (*WorkspaceV2Event) isEvent() {}

// CreateWorkspaceEvent is emitted when a workspace is created.
//
// Wire: createworkspace>>3
type CreateWorkspaceEvent struct {
	WorkspaceName string
}

// Name returns EventCreateWorkspace.
func (*CreateWorkspaceEvent) Name() string { return EventCreateWorkspace }

func (*CreateWorkspaceEvent) isEvent() {}

// CreateWorkspaceV2Event is the v2 variant, providing the workspace ID as well.
//
// Wire: createworkspacev2>>3,3
type CreateWorkspaceV2Event struct {
	WorkspaceID   int
	WorkspaceName string
}

// Name returns EventCreateWorkspaceV2.
func (*CreateWorkspaceV2Event) Name() string { return EventCreateWorkspaceV2 }

func (*CreateWorkspaceV2Event) isEvent() {}

// DestroyWorkspaceEvent is emitted when a workspace is destroyed.
//
// Wire: destroyworkspace>>3
type DestroyWorkspaceEvent struct {
	WorkspaceName string
}

// Name returns EventDestroyWorkspace.
func (*DestroyWorkspaceEvent) Name() string { return EventDestroyWorkspace }

func (*DestroyWorkspaceEvent) isEvent() {}

// DestroyWorkspaceV2Event is the v2 variant, providing the workspace ID as well.
//
// Wire: destroyworkspacev2>>3,3
type DestroyWorkspaceV2Event struct {
	WorkspaceID   int
	WorkspaceName string
}

// Name returns EventDestroyWorkspaceV2.
func (*DestroyWorkspaceV2Event) Name() string { return EventDestroyWorkspaceV2 }

func (*DestroyWorkspaceV2Event) isEvent() {}

// OpenLayerEvent is emitted when a layer surface is mapped.
//
// Wire: openlayer>>waybar
type OpenLayerEvent struct {
	Namespace string
}

// Name returns EventOpenLayer.
func (*OpenLayerEvent) Name() string { return EventOpenLayer }

func (*OpenLayerEvent) isEvent() {}

// CloseLayerEvent is emitted when a layer surface is unmapped.
//
// Wire: closelayer>>waybar
type CloseLayerEvent struct {
	Namespace string
}

// Name returns EventCloseLayer.
func (*CloseLayerEvent) Name() string { return EventCloseLayer }

func (*CloseLayerEvent) isEvent() {}

// WindowTitleEvent is emitted when a window title changes.
//
// Wire: windowtitle>>56157f7c25a0
type WindowTitleEvent struct {
	Address WindowAddress
}

// Name returns EventWindowTitle.
func (*WindowTitleEvent) Name() string { return EventWindowTitle }

func (*WindowTitleEvent) isEvent() {}

// WindowTitleV2Event is the v2 variant, providing the new title as well. The
// title may contain commas.
//
// Wire: windowtitlev2>>56157f7c25a0,zsh
type WindowTitleV2Event struct {
	Address WindowAddress
	Title   string
}

// Name returns EventWindowTitleV2.
func (*WindowTitleV2Event) Name() string { return EventWindowTitleV2 }

func (*WindowTitleV2Event) isEvent() {}

// OpenWindowEvent is emitted when a window is opened. The title may contain
// commas.
//
// Wire: openwindow>>56157e75fd30,2,firefox,Mozilla Firefox
type OpenWindowEvent struct {
	Address       WindowAddress
	WorkspaceName string
	Class         string
	Title         string
}

// Name returns EventOpenWindow.
func (*OpenWindowEvent) Name() string { return EventOpenWindow }

func (*OpenWindowEvent) isEvent() {}

// CloseWindowEvent is emitted when a window is closed.
//
// Wire: closewindow>>56157f7c25a0
type CloseWindowEvent struct {
	Address WindowAddress
}

// Name returns EventCloseWindow.
func (*CloseWindowEvent) Name() string { return EventCloseWindow }

func (*CloseWindowEvent) isEvent() {}

// FocusedMonEvent is emitted when the active monitor changes.
//
// Wire: focusedmon>>DP-1,2
type FocusedMonEvent struct {
	MonitorName   string
	WorkspaceName string
}

// Name returns EventFocusedMon.
func (*FocusedMonEvent) Name() string { return EventFocusedMon }

func (*FocusedMonEvent) isEvent() {}

// FocusedMonV2Event is the v2 variant, providing the workspace ID instead of
// its name.
//
// Wire: focusedmonv2>>DP-1,2
type FocusedMonV2Event struct {
	MonitorName string
	WorkspaceID int
}

// Name returns EventFocusedMonV2.
func (*FocusedMonV2Event) Name() string { return EventFocusedMonV2 }

func (*FocusedMonV2Event) isEvent() {}

// FullscreenEvent is emitted when a window's fullscreen status changes.
//
// Wire: fullscreen>>1
type FullscreenEvent struct {
	Fullscreen bool
}

// Name returns EventFullscreen.
func (*FullscreenEvent) Name() string { return EventFullscreen }

func (*FullscreenEvent) isEvent() {}

// MonitorRemovedEvent is emitted when a monitor is disconnected.
//
// Wire: monitorremoved>>DP-1
type MonitorRemovedEvent struct {
	MonitorName string
}

// Name returns EventMonitorRemoved.
func (*MonitorRemovedEvent) Name() string { return EventMonitorRemoved }

func (*MonitorRemovedEvent) isEvent() {}

// MonitorRemovedV2Event is the v2 variant, providing the monitor ID and
// description as well. The description may contain commas.
//
// Wire: monitorremovedv2>>0,DP-1,Dell U2720Q
type MonitorRemovedV2Event struct {
	MonitorID          int
	MonitorName        string
	MonitorDescription string
}

// Name returns EventMonitorRemovedV2.
func (*MonitorRemovedV2Event) Name() string { return EventMonitorRemovedV2 }

func (*MonitorRemovedV2Event) isEvent() {}

// MonitorAddedEvent is emitted when a monitor is connected.
//
// Wire: monitoradded>>DP-1
type MonitorAddedEvent struct {
	MonitorName string
}

// Name returns EventMonitorAdded.
func (*MonitorAddedEvent) Name() string { return EventMonitorAdded }

func (*MonitorAddedEvent) isEvent() {}

// MonitorAddedV2Event is the v2 variant, providing the monitor ID and
// description as well. The description may contain commas.
//
// Wire: monitoraddedv2>>0,DP-1,Dell U2720Q
type MonitorAddedV2Event struct {
	MonitorID          int
	MonitorName        string
	MonitorDescription string
}

// Name returns EventMonitorAddedV2.
func (*MonitorAddedV2Event) Name() string { return EventMonitorAddedV2 }

func (*MonitorAddedV2Event) isEvent() {}

// KillEvent is emitted when a window is killed via hyprctl kill.
//
// Wire: kill>>56157f7c25a0
type KillEvent struct {
	Address WindowAddress
}

// Name returns EventKill.
func (*KillEvent) Name() string { return EventKill }

func (*KillEvent) isEvent() {}

// MoveWindowEvent is emitted when a window is moved to a different workspace.
//
// Wire: movewindow>>56157f7c25a0,2
type MoveWindowEvent struct {
	Address       WindowAddress
	WorkspaceName string
}

// Name returns EventMoveWindow.
func (*MoveWindowEvent) Name() string { return EventMoveWindow }

func (*MoveWindowEvent) isEvent() {}

// MoveWindowV2Event is the v2 variant, providing the workspace ID as well.
//
// Wire: movewindowv2>>56157f7c25a0,2,web
type MoveWindowV2Event struct {
	Address       WindowAddress
	WorkspaceID   int
	WorkspaceName string
}

// Name returns EventMoveWindowV2.
func (*MoveWindowV2Event) Name() string { return EventMoveWindowV2 }

func (*MoveWindowV2Event) isEvent() {}

// MoveWorkspaceEvent is emitted when a workspace is moved to a different
// monitor.
//
// Wire: moveworkspace>>2,DP-1
type MoveWorkspaceEvent struct {
	WorkspaceName string
	MonitorName   string
}

// Name returns EventMoveWorkspace.
func (*MoveWorkspaceEvent) Name() string { return EventMoveWorkspace }

func (*MoveWorkspaceEvent) isEvent() {}

// MoveWorkspaceV2Event is the v2 variant, providing the workspace ID as well.
//
// Wire: moveworkspacev2>>2,web,DP-1
type MoveWorkspaceV2Event struct {
	WorkspaceID   int
	WorkspaceName string
	MonitorName   string
}

// Name returns EventMoveWorkspaceV2.
func (*MoveWorkspaceV2Event) Name() string { return EventMoveWorkspaceV2 }

func (*MoveWorkspaceV2Event) isEvent() {}

// RenameWorkspaceEvent is emitted when a workspace is renamed. The new name
// may contain commas.
//
// Wire: renameworkspace>>2,web
type RenameWorkspaceEvent struct {
	WorkspaceID int
	NewName     string
}

// Name returns EventRenameWorkspace.
func (*RenameWorkspaceEvent) Name() string { return EventRenameWorkspace }

func (*RenameWorkspaceEvent) isEvent() {}

// ChangeWorkspaceIDEvent is emitted when a workspace's ID changes.
//
// Wire: changeworkspaceid>>6,77
type ChangeWorkspaceIDEvent struct {
	OldID int
	NewID int
}

// Name returns EventChangeWorkspaceID.
func (*ChangeWorkspaceIDEvent) Name() string { return EventChangeWorkspaceID }

func (*ChangeWorkspaceIDEvent) isEvent() {}

// ActiveSpecialEvent is emitted when the special workspace opened on a monitor
// changes. WorkspaceName is empty when the special workspace closed.
//
// Wire: activespecial>>special:magic,DP-1
type ActiveSpecialEvent struct {
	WorkspaceName string
	MonitorName   string
}

// Name returns EventActiveSpecial.
func (*ActiveSpecialEvent) Name() string { return EventActiveSpecial }

func (*ActiveSpecialEvent) isEvent() {}

// ActiveSpecialV2Event is the v2 variant, providing the workspace ID as well.
// WorkspaceID is 0 and WorkspaceName empty when the special workspace closed.
//
// Wire: activespecialv2>>-98,special:magic,DP-1
type ActiveSpecialV2Event struct {
	WorkspaceID   int
	WorkspaceName string
	MonitorName   string
}

// Name returns EventActiveSpecialV2.
func (*ActiveSpecialV2Event) Name() string { return EventActiveSpecialV2 }

func (*ActiveSpecialV2Event) isEvent() {}

// ActiveLayoutEvent is emitted on a layout change of the active keyboard.
//
// Wire: activelayout>>at-translated-set-2-keyboard,English (US)
type ActiveLayoutEvent struct {
	KeyboardName string
	LayoutName   string
}

// Name returns EventActiveLayout.
func (*ActiveLayoutEvent) Name() string { return EventActiveLayout }

func (*ActiveLayoutEvent) isEvent() {}

// SubmapEvent is emitted when a keybind submap changes. An empty Submap means
// the default submap.
//
// Wire: submap>>resize
type SubmapEvent struct {
	Submap string
}

// Name returns EventSubmap.
func (*SubmapEvent) Name() string { return EventSubmap }

func (*SubmapEvent) isEvent() {}

// ChangeFloatingModeEvent is emitted when a window changes its floating mode.
//
// Wire: changefloatingmode>>56157f7c25a0,1
type ChangeFloatingModeEvent struct {
	Address  WindowAddress
	Floating bool
}

// Name returns EventChangeFloatingMode.
func (*ChangeFloatingModeEvent) Name() string { return EventChangeFloatingMode }

func (*ChangeFloatingModeEvent) isEvent() {}

// UrgentEvent is emitted when a window requests an urgent state.
//
// Wire: urgent>>56157f7c25a0
type UrgentEvent struct {
	Address WindowAddress
}

// Name returns EventUrgent.
func (*UrgentEvent) Name() string { return EventUrgent }

func (*UrgentEvent) isEvent() {}

// ScreenCastEvent is emitted when the screencopy state of a client changes.
//
// Wire: screencast>>1,0
type ScreenCastEvent struct {
	Sharing bool
	Owner   ScreenCastOwner
}

// Name returns EventScreenCast.
func (*ScreenCastEvent) Name() string { return EventScreenCast }

func (*ScreenCastEvent) isEvent() {}

// ScreenCastV2Event is the v2 variant, providing the identifier of the shared
// target as well.
//
// Wire: screencastv2>>1,1,56157f7c25a0
type ScreenCastV2Event struct {
	Sharing    bool
	Owner      ScreenCastOwner
	TargetName string
}

// Name returns EventScreenCastV2.
func (*ScreenCastV2Event) Name() string { return EventScreenCastV2 }

func (*ScreenCastV2Event) isEvent() {}

// ToggleGroupEvent is emitted when the togglegroup dispatcher is used. Grouped
// false means the group was destroyed; Addresses lists the windows involved.
//
// Wire: togglegroup>>0,64cea2525760,64cea2522380
type ToggleGroupEvent struct {
	Grouped   bool
	Addresses []WindowAddress
}

// Name returns EventToggleGroup.
func (*ToggleGroupEvent) Name() string { return EventToggleGroup }

func (*ToggleGroupEvent) isEvent() {}

// MoveIntoGroupEvent is emitted when a window is merged into a group.
//
// Wire: moveintogroup>>56157f7c25a0
type MoveIntoGroupEvent struct {
	Address WindowAddress
}

// Name returns EventMoveIntoGroup.
func (*MoveIntoGroupEvent) Name() string { return EventMoveIntoGroup }

func (*MoveIntoGroupEvent) isEvent() {}

// MoveOutOfGroupEvent is emitted when a window is removed from a group.
//
// Wire: moveoutofgroup>>56157f7c25a0
type MoveOutOfGroupEvent struct {
	Address WindowAddress
}

// Name returns EventMoveOutOfGroup.
func (*MoveOutOfGroupEvent) Name() string { return EventMoveOutOfGroup }

func (*MoveOutOfGroupEvent) isEvent() {}

// IgnoreGroupLockEvent is emitted when ignoregrouplock is toggled.
//
// Wire: ignoregrouplock>>1
type IgnoreGroupLockEvent struct {
	Ignored bool
}

// Name returns EventIgnoreGroupLock.
func (*IgnoreGroupLockEvent) Name() string { return EventIgnoreGroupLock }

func (*IgnoreGroupLockEvent) isEvent() {}

// LockGroupsEvent is emitted when lockgroups is toggled.
//
// Wire: lockgroups>>1
type LockGroupsEvent struct {
	Locked bool
}

// Name returns EventLockGroups.
func (*LockGroupsEvent) Name() string { return EventLockGroups }

func (*LockGroupsEvent) isEvent() {}

// ConfigReloadedEvent is emitted when the config is done reloading.
//
// Wire: configreloaded>>
type ConfigReloadedEvent struct{}

// Name returns EventConfigReloaded.
func (*ConfigReloadedEvent) Name() string { return EventConfigReloaded }

func (*ConfigReloadedEvent) isEvent() {}

// PinEvent is emitted when a window is pinned or unpinned.
//
// Wire: pin>>56157f7c25a0,1
type PinEvent struct {
	Address WindowAddress
	Pinned  bool
}

// Name returns EventPin.
func (*PinEvent) Name() string { return EventPin }

func (*PinEvent) isEvent() {}

// MinimizedEvent is emitted when an external taskbar-like app requests a window
// to be minimized.
//
// Wire: minimized>>56157f7c25a0,1
type MinimizedEvent struct {
	Address   WindowAddress
	Minimized bool
}

// Name returns EventMinimized.
func (*MinimizedEvent) Name() string { return EventMinimized }

func (*MinimizedEvent) isEvent() {}

// BellEvent is emitted when an app requests to ring the system bell via
// xdg-system-bell-v1. Address may be empty.
//
// Wire: bell>>56157f7c25a0
type BellEvent struct {
	Address WindowAddress
}

// Name returns EventBell.
func (*BellEvent) Name() string { return EventBell }

func (*BellEvent) isEvent() {}

// CustomEvent is emitted by the event dispatcher. Data is whatever string the
// caller passed, not split.
//
// Wire: custom>>hello, world
type CustomEvent struct {
	Data string
}

// Name returns EventCustom.
func (*CustomEvent) Name() string { return EventCustom }

func (*CustomEvent) isEvent() {}
