package hypr

import (
	"context"
	"reflect"
	"testing"
)

func TestMonitors(t *testing.T) {
	t.Parallel()

	monitors, err := Monitors(context.Background(), fixtureOpts(t, "monitors.json")...)
	if err != nil {
		t.Fatalf("Monitors: %v", err)
	}
	if len(monitors) != 2 {
		t.Fatalf("got %d monitors, want 2", len(monitors))
	}

	want := Monitor{
		ID: 0, Name: "eDP-1", Description: "Lenovo Group Limited LEN140WUXGA",
		Make: "Lenovo Group Limited", Model: "LEN140WUXGA  ",
		Width: 1440, Height: 900, PhysicalWidth: 300, PhysicalHeight: 190, RefreshRate: 60.003,
		ActiveWorkspace: WorkspaceRef{ID: 3, Name: "3"}, SpecialWorkspace: WorkspaceRef{},
		Reserved: [4]int{0, 32, 0, 0}, Scale: 1, Focused: true, DPMSStatus: true,
		Solitary: "0", SolitaryBlockedBy: []string{"WINDOWED", "CANDIDATE"},
		TearingBlockedBy: []string{"NOT_TORN", "USER", "CANDIDATE", "HW_CURSOR"},
		DirectScanoutTo:  "0", DirectScanoutBlockedBy: []string{"USER", "CANDIDATE"},
		CurrentFormat: "XRGB8888", MirrorOf: "none",
		AvailableModes: []string{
			"1920x1200@60.00Hz", "1920x1080@60.00Hz", "1600x1200@60.00Hz", "1680x1050@60.00Hz",
			"1280x1024@60.00Hz", "1440x900@60.00Hz", "1280x800@60.00Hz", "1280x720@60.00Hz",
			"1024x768@60.00Hz", "800x600@60.00Hz", "640x480@60.00Hz",
		},
		ColorManagementPreset: "srgb", SDRBrightness: 1, SDRSaturation: 1, SDRMinLuminance: 0.2, SDRMaxLuminance: 80,
		HardwareCursorsInUse: true,
	}
	if !reflect.DeepEqual(monitors[0], want) {
		t.Errorf("monitors[0]\n got %+v\nwant %+v", monitors[0], want)
	}

	if m := monitors[1]; m.ID != 1 || m.Name != "DP-2" || m.ActiveWorkspace != (WorkspaceRef{ID: 5, Name: "5"}) || m.RefreshRate != 29.985 || m.Serial != "SERIAL" {
		t.Errorf("monitors[1] = %+v", m)
	}
}
