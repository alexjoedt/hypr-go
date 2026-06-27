package hypr

import (
	"encoding/json"
	"strings"
)

// WindowAddress identifies a window the way the event socket emits it: hex
// without the 0x prefix. hyprctl -j adds the prefix; UnmarshalJSON strips it,
// so both sides compare equal.
type WindowAddress string

// String returns the address as the event socket emits it, without 0x.
func (a WindowAddress) String() string {
	return string(a)
}

// UnmarshalJSON accepts a JSON string and strips a leading 0x.
func (a *WindowAddress) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*a = WindowAddress(strings.TrimPrefix(s, "0x"))

	return nil
}
