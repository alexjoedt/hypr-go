// Package hypr talks to a running Hyprland instance over its two sockets.
//
// Commands are package-level functions: [Dispatch], [Request] and [Batch] dial
// the command socket, send, read the reply and close. Events are one stateful
// type: [Listen] dials the event socket and returns a [Listener] that owns it
// until Close; [On] registers typed handlers and [Listener.Events] iterates.
//
// A reply Hyprland did not answer with ok is a [*ReplyError]; an event this
// version has no type for is an [*UnknownEvent]. The escape hatches for a newer
// Hyprland are [RawCommand] and [Request] on the command side and
// [UnknownEvent] on the event side.
//
// # Finding a name
//
// The package is flat on purpose. Every identifier is derived from the name
// Hyprland uses, so the wiki name is what to type:
//
//   - A dispatcher is its Lua path in CamelCase: hl.dsp.window.close is
//     [WindowClose], hl.dsp.focus is [Focus]. An overloaded dispatcher has one
//     struct per mode, the mode appended: [FocusWindow], [FocusMonitor].
//   - An event is its name plus Event: workspacev2 is [WorkspaceV2Event].
//   - A query is the hyprctl name as a function. A list reply returns the
//     singular noun, [Clients] returns [Client]. An object reply returns the
//     name plus Info, [Version] returns [VersionInfo].
//   - A selector is the noun, By and the key: [WindowByClass],
//     [WorkspaceByID], [MonitorByName]. Fixed values are constants on the
//     noun: [WindowActive], [WorkspaceNext], [MonitorCurrent].
//   - An enum value carries its type as prefix: [ActionToggle],
//     [DirectionLeft], [FullscreenStateNone], [CornerTopLeft].
//
// Hyprland says client in hyprctl -j clients and window everywhere else, so
// the data struct is [Client] and the selectors and dispatchers say Window.
// The one dispatcher not named after its path is [EmitEvent] for
// hl.dsp.event, because Event is the interface.
package hypr
