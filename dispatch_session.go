package hypr

import (
	"strconv"
	"time"
)

// KeyState is the state argument of SendKeyState.
type KeyState int

// The KeyState values, in Hyprland's order.
const (
	KeyDown KeyState = iota
	KeyUp
	KeyRepeat
)

// String returns the Lua value: "down", "up" or "repeat".
func (s KeyState) String() string {
	switch s {
	case KeyUp:
		return "up"
	case KeyRepeat:
		return "repeat"
	default:
		return "down"
	}
}

// ExecCmd wraps hl.dsp.exec_cmd: it runs cmd through sh -c. The string is
// passed through verbatim; quoting is the caller's job.
func ExecCmd(cmd string) Command {
	return RawCommand(positional("exec_cmd", quote(cmd)))
}

// ExecRaw wraps hl.dsp.exec_raw: it runs cmd without a shell.
func ExecRaw(cmd string) Command {
	return RawCommand(positional("exec_raw", quote(cmd)))
}

// Exit wraps hl.dsp.exit: it ends the Hyprland session.
func Exit() Command {
	return RawCommand(positional("exit"))
}

// ReloadConfig wraps hl.dsp.reload_config. Hyprland 0.56.2 does not register
// it and replies with an error; on that version use Request with "/reload"
// instead.
func ReloadConfig() Command {
	return RawCommand(positional("reload_config"))
}

// Submap wraps hl.dsp.submap: it enters the named keybind submap; "reset" or
// the empty string leaves it.
func Submap(name string) Command {
	return RawCommand(positional("submap", quote(name)))
}

// Pass wraps hl.dsp.pass: it passes the triggering key to the window instead
// of handling it.
type Pass struct {
	Window WindowSelector
}

// Command implements Command.
func (c Pass) Command() string {
	return dispatch("pass", str("window", string(c.Window)))
}

// SendShortcut wraps hl.dsp.send_shortcut: it sends a key press with the
// modifiers, "SUPER SHIFT" or "" for none, to the window. Key is an xkb keysym
// name, "code:N" or "mouse:N". Window defaults to the focused one.
type SendShortcut struct {
	Mods   string
	Key    string
	Window WindowSelector
}

// Command implements Command.
func (c SendShortcut) Command() string {
	return dispatch("send_shortcut",
		str("mods", c.Mods),
		str("key", c.Key),
		optStr("window", string(c.Window)),
	)
}

// SendKeyState wraps hl.dsp.send_key_state: like SendShortcut, but with an
// explicit key state. Window defaults to the focused one.
type SendKeyState struct {
	Mods   string
	Key    string
	State  KeyState
	Window WindowSelector
}

// Command implements Command.
func (c SendKeyState) Command() string {
	return dispatch("send_key_state",
		str("mods", c.Mods),
		str("key", c.Key),
		str("state", c.State.String()),
		optStr("window", string(c.Window)),
	)
}

// Layout wraps hl.dsp.layout: it sends a message to the active layout,
// "togglesplit" for dwindle.
func Layout(message string) Command {
	return RawCommand(positional("layout", quote(message)))
}

// DPMS wraps hl.dsp.dpms: it toggles, enables or disables the display of a
// monitor. Monitor defaults to every monitor.
type DPMS struct {
	Action  Action
	Monitor MonitorSelector
}

// Command implements Command.
func (c DPMS) Command() string {
	return dispatch("dpms",
		optAction("action", c.Action),
		optStr("monitor", string(c.Monitor)),
	)
}

// EmitEvent wraps hl.dsp.event: it emits "custom>>data" on the event socket,
// which Parse turns into a CustomEvent.
func EmitEvent(data string) Command {
	return RawCommand(positional("event", quote(data)))
}

// Global wraps hl.dsp.global: it triggers a global shortcut registered by an
// app, "appid:name".
func Global(name string) Command {
	return RawCommand(positional("global", quote(name)))
}

// ForceRendererReload wraps hl.dsp.force_renderer_reload.
func ForceRendererReload() Command {
	return RawCommand(positional("force_renderer_reload"))
}

// ForceIdle wraps hl.dsp.force_idle: it reports the session idle for that
// long. Sub-second durations are sent as fractional seconds.
func ForceIdle(d time.Duration) Command {
	return RawCommand(positional("force_idle", strconv.FormatFloat(d.Seconds(), 'f', -1, 64)))
}

// ReleaseInputCapture wraps hl.dsp.release_input_capture: it releases an
// input capture held by a client.
func ReleaseInputCapture() Command {
	return RawCommand(positional("release_input_capture"))
}

// NoOp wraps hl.dsp.no_op: it does nothing and replies ok.
func NoOp() Command {
	return RawCommand(positional("no_op"))
}
