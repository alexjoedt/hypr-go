package hypr

import (
	"encoding/json"
	"testing"
)

func TestWindowAddress(t *testing.T) {
	t.Parallel()

	if got := WindowAddress("5a1b2c3d").String(); got != "5a1b2c3d" {
		t.Errorf("String() = %q, want 5a1b2c3d", got)
	}

	tests := []struct {
		name string
		json string
		want WindowAddress
	}{
		{name: "hyprctl prefix", json: `"0x5a1b2c3d"`, want: "5a1b2c3d"},
		{name: "no prefix", json: `"5a1b2c3d"`, want: "5a1b2c3d"},
		{name: "empty", json: `""`, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got struct {
				Address WindowAddress `json:"address"`
			}
			if err := json.Unmarshal([]byte(`{"address":`+tt.json+`}`), &got); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			if got.Address != tt.want {
				t.Errorf("Address = %q, want %q", got.Address, tt.want)
			}
		})
	}

	t.Run("not a string", func(t *testing.T) {
		t.Parallel()

		var a WindowAddress
		if err := json.Unmarshal([]byte(`123`), &a); err == nil {
			t.Error("Unmarshal accepted a number")
		}
	})
}
