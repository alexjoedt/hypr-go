package hypr

import "testing"

func TestDispatchSerializer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		path   string
		fields []field
		want   string
	}{
		{name: "no fields", path: "no_op", want: "hl.dsp.no_op()"},
		{name: "all omitted", path: "window.close", fields: []field{optStr("window", "")}, want: "hl.dsp.window.close()"},
		{name: "string", path: "window.close", fields: []field{str("window", "class:^kitty$")}, want: `hl.dsp.window.close({ window = "class:^kitty$" })`},
		{name: "empty required string", path: "workspace.rename", fields: []field{str("name", "")}, want: `hl.dsp.workspace.rename({ name = "" })`},
		{name: "int", path: "window.signal", fields: []field{num("signal", 15)}, want: "hl.dsp.window.signal({ signal = 15 })"},
		{name: "zero required int", path: "cursor.move", fields: []field{num("x", 0), num("y", -20)}, want: "hl.dsp.cursor.move({ x = 0, y = -20 })"},
		{name: "optional int omits zero", path: "group.active", fields: []field{optNum("index", 0), optNum("other", 2)}, want: "hl.dsp.group.active({ other = 2 })"},
		{name: "bool", path: "group.move_window", fields: []field{boolean("forward", false)}, want: "hl.dsp.group.move_window({ forward = false })"},
		{name: "optional bool omits false", path: "window.move", fields: []field{optBool("a", false), optBool("b", true)}, want: "hl.dsp.window.move({ b = true })"},
		{name: "action toggle omitted", path: "window.pin", fields: []field{optAction("action", ActionToggle)}, want: "hl.dsp.window.pin()"},
		{name: "action on", path: "window.pin", fields: []field{optAction("action", ActionEnable)}, want: `hl.dsp.window.pin({ action = "on" })`},
		{name: "action off", path: "window.pin", fields: []field{optAction("action", ActionDisable)}, want: `hl.dsp.window.pin({ action = "off" })`},
		{
			name: "order is the order given, omitted fields leave no gap",
			path: "window.fullscreen",
			fields: []field{
				optStr("mode", "maximized"),
				optStr("action", ""),
				optBool("layout_aware", false),
				optStr("window", "active"),
			},
			want: `hl.dsp.window.fullscreen({ mode = "maximized", window = "active" })`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := dispatch(tt.path, tt.fields...); got != tt.want {
				t.Errorf("dispatch() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPositional(t *testing.T) {
	t.Parallel()

	if got := positional("exit"); got != "hl.dsp.exit()" {
		t.Errorf("positional() = %q", got)
	}
	if got := positional("submap", quote("resize")); got != `hl.dsp.submap("resize")` {
		t.Errorf("positional() = %q", got)
	}
	if got := positional("exec_cmd", quote("kitty"), "{ float = true }"); got != `hl.dsp.exec_cmd("kitty", { float = true })` {
		t.Errorf("positional() = %q", got)
	}
}

func TestQuote(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"", `""`},
		{"kitty", `"kitty"`},
		{`class:^(kitty)$`, `"class:^(kitty)$"`},
		{`say "hi"`, `"say \"hi\""`},
		{`back\slash`, `"back\\slash"`},
		{"line\nbreak", `"line\nbreak"`},
		{"tab\there", `"tab\there"`},
		{"cr\rhere", `"cr\rhere"`},
		{"bell\x07", `"bell\7"`},
		{"del\x7f", `"del\127"`},
		{"ünïcödé", `"ünïcödé"`},
		{"sh -c 'echo $HOME'", `"sh -c 'echo $HOME'"`},
	}

	for _, tt := range tests {
		if got := quote(tt.in); got != tt.want {
			t.Errorf("quote(%q) = %s, want %s", tt.in, got, tt.want)
		}
	}
}

func TestActionString(t *testing.T) {
	t.Parallel()

	tests := map[Action]string{ActionToggle: "toggle", ActionEnable: "on", ActionDisable: "off", Action(7): "toggle"}
	for a, want := range tests {
		if got := a.String(); got != want {
			t.Errorf("Action(%d).String() = %q, want %q", int(a), got, want)
		}
	}
}
