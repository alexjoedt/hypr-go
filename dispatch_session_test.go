package hypr

import (
	"reflect"
	"testing"
	"time"
)

var sessionCommandTests = []commandTest{
	{ExecCmd("kitty"), `hl.dsp.exec_cmd("kitty")`},
	{ExecCmd(`sh -c "echo \"hi\""`), `hl.dsp.exec_cmd("sh -c \"echo \\\"hi\\\"\"")`},
	{ExecRaw("kitty --single-instance"), `hl.dsp.exec_raw("kitty --single-instance")`},
	{Exit(), `hl.dsp.exit()`},
	{ReloadConfig(), `hl.dsp.reload_config()`},
	{Submap("resize"), `hl.dsp.submap("resize")`},
	{Submap(""), `hl.dsp.submap("")`},

	{Pass{Window: WindowByClass("^obs$")}, `hl.dsp.pass({ window = "class:^obs$" })`},
	{Pass{}, `hl.dsp.pass({ window = "" })`},

	{SendShortcut{Mods: "SUPER", Key: "a"}, `hl.dsp.send_shortcut({ mods = "SUPER", key = "a" })`},
	{SendShortcut{Mods: "", Key: "code:38", Window: kitty}, `hl.dsp.send_shortcut({ mods = "", key = "code:38", window = "class:^kitty$" })`},

	{SendKeyState{Mods: "", Key: "a", State: KeyDown}, `hl.dsp.send_key_state({ mods = "", key = "a", state = "down" })`},
	{SendKeyState{Mods: "SUPER SHIFT", Key: "q", State: KeyUp, Window: kitty}, `hl.dsp.send_key_state({ mods = "SUPER SHIFT", key = "q", state = "up", window = "class:^kitty$" })`},
	{SendKeyState{Key: "a", State: KeyRepeat}, `hl.dsp.send_key_state({ mods = "", key = "a", state = "repeat" })`},

	{Layout("togglesplit"), `hl.dsp.layout("togglesplit")`},

	{DPMS{}, `hl.dsp.dpms()`},
	{DPMS{Action: ActionDisable, Monitor: MonitorByName("DP-1")}, `hl.dsp.dpms({ action = "off", monitor = "DP-1" })`},
	{DPMS{Action: ActionEnable}, `hl.dsp.dpms({ action = "on" })`},

	{EmitEvent("hello"), `hl.dsp.event("hello")`},
	{Global("com.example.app:myshortcut"), `hl.dsp.global("com.example.app:myshortcut")`},
	{ForceRendererReload(), `hl.dsp.force_renderer_reload()`},
	{ForceIdle(60 * time.Second), `hl.dsp.force_idle(60)`},
	{ForceIdle(1500 * time.Millisecond), `hl.dsp.force_idle(1.5)`},
	{ReleaseInputCapture(), `hl.dsp.release_input_capture()`},
	{NoOp(), `hl.dsp.no_op()`},
}

func TestSessionDispatchers(t *testing.T) {
	t.Parallel()

	runCommandTests(t, sessionCommandTests)
}

// TestEmitEventRoundTrip pairs EmitEvent with Parse: the payload the
// dispatcher carries is the payload the custom event delivers.
func TestEmitEventRoundTrip(t *testing.T) {
	t.Parallel()

	for _, data := range []string{"hello", "hello, world", "a>>b", ""} {
		if got := EmitEvent(data).Command(); got != `hl.dsp.event(`+quote(data)+`)` {
			t.Errorf("EmitEvent(%q) = %s", data, got)
		}

		got, err := Parse([]byte(EventCustom + ">>" + data))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if want := (&CustomEvent{Data: data}); !reflect.DeepEqual(got, want) {
			t.Errorf("Parse(custom>>%s) = %+v, want %+v", data, got, want)
		}
	}
}

func TestKeyStateString(t *testing.T) {
	t.Parallel()

	tests := map[KeyState]string{KeyDown: "down", KeyUp: "up", KeyRepeat: "repeat", KeyState(9): "down"}
	for s, want := range tests {
		if got := s.String(); got != want {
			t.Errorf("KeyState(%d).String() = %q, want %q", int(s), got, want)
		}
	}
}
