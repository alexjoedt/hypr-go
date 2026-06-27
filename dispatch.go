package hypr

import (
	"strconv"
	"strings"
)

// Action is the toggle argument many dispatchers take. The zero value is
// Hyprland's default, toggle, and is omitted from the wire string.
type Action int

// The Action values, in Hyprland's order.
const (
	ActionToggle Action = iota
	ActionEnable
	ActionDisable
)

// String returns the Lua value: "toggle", "on" or "off".
func (a Action) String() string {
	switch a {
	case ActionEnable:
		return "on"
	case ActionDisable:
		return "off"
	default:
		return "toggle"
	}
}

// field is one key of a dispatcher's Lua table. A field with omit set is
// left out, which is how zero values fall back to Hyprland's default.
type field struct {
	key   string
	value string
	omit  bool
}

// str always emits a quoted string; optStr omits an empty one.
func str(key, v string) field        { return field{key: key, value: quote(v)} }
func optStr(key, v string) field     { return field{key: key, value: quote(v), omit: v == ""} }
func num(key string, v int) field    { return field{key: key, value: strconv.Itoa(v)} }
func optNum(key string, v int) field { return field{key: key, value: strconv.Itoa(v), omit: v == 0} }

// boolean always emits; optBool omits false.
func boolean(key string, v bool) field { return field{key: key, value: strconv.FormatBool(v)} }
func optBool(key string, v bool) field {
	return field{key: key, value: strconv.FormatBool(v), omit: !v}
}

// optFalse emits key = false only when v is set. It carries a negated Go field
// (NoFollow) onto a Lua key whose default is true (follow).
func optFalse(key string, v bool) field { return field{key: key, value: "false", omit: !v} }

// optEnum emits a quoted enum value unless it is the default.
func optEnum(key, value string, isDefault bool) field {
	return field{key: key, value: quote(value), omit: isDefault}
}

// optAction omits the default, toggle.
func optAction(key string, a Action) field {
	return field{key: key, value: quote(a.String()), omit: a == ActionToggle}
}

// dispatch builds the wire expression for a table dispatcher: the Lua path
// after hl.dsp. plus the fields that are not omitted, in the order given.
// With nothing to emit the call carries no table at all.
func dispatch(path string, fields ...field) string {
	var b strings.Builder
	b.WriteString("hl.dsp.")
	b.WriteString(path)
	b.WriteByte('(')

	n := 0
	for _, f := range fields {
		if f.omit {
			continue
		}
		if n == 0 {
			b.WriteString("{ ")
		} else {
			b.WriteString(", ")
		}
		b.WriteString(f.key)
		b.WriteString(" = ")
		b.WriteString(f.value)
		n++
	}
	if n > 0 {
		b.WriteString(" }")
	}

	b.WriteByte(')')

	return b.String()
}

// positional builds the wire expression for a dispatcher that takes its
// arguments in order rather than as a table: hl.dsp.submap("resize").
func positional(path string, values ...string) string {
	return "hl.dsp." + path + "(" + strings.Join(values, ", ") + ")"
}

// quote writes s as a double-quoted Lua string literal. Backslash, the quote
// and control characters are escaped, so the value can carry any regex or
// shell string verbatim.
func quote(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\\', '"':
			b.WriteByte('\\')
			b.WriteByte(c)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if c < 0x20 || c == 0x7f {
				b.WriteString(`\`)
				b.WriteString(strconv.Itoa(int(c)))
			} else {
				b.WriteByte(c)
			}
		}
	}
	b.WriteByte('"')

	return b.String()
}
