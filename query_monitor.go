package hypr

import "context"

// Monitor is one output as hyprctl -j monitors reports it.
type Monitor struct {
	ID                     int          `json:"id"`
	Name                   string       `json:"name"`
	Description            string       `json:"description"`
	Make                   string       `json:"make"`
	Model                  string       `json:"model"`
	Serial                 string       `json:"serial"`
	Width                  int          `json:"width"`
	Height                 int          `json:"height"`
	PhysicalWidth          int          `json:"physicalWidth"`
	PhysicalHeight         int          `json:"physicalHeight"`
	RefreshRate            float64      `json:"refreshRate"`
	X                      int          `json:"x"`
	Y                      int          `json:"y"`
	ActiveWorkspace        WorkspaceRef `json:"activeWorkspace"`
	SpecialWorkspace       WorkspaceRef `json:"specialWorkspace"`
	Reserved               [4]int       `json:"reserved"`
	Scale                  float64      `json:"scale"`
	Transform              int          `json:"transform"`
	Focused                bool         `json:"focused"`
	DPMSStatus             bool         `json:"dpmsStatus"`
	VRR                    bool         `json:"vrr"`
	Solitary               string       `json:"solitary"`
	SolitaryBlockedBy      []string     `json:"solitaryBlockedBy"`
	ActivelyTearing        bool         `json:"activelyTearing"`
	TearingBlockedBy       []string     `json:"tearingBlockedBy"`
	DirectScanoutTo        string       `json:"directScanoutTo"`
	DirectScanoutBlockedBy []string     `json:"directScanoutBlockedBy"`
	Disabled               bool         `json:"disabled"`
	CurrentFormat          string       `json:"currentFormat"`
	MirrorOf               string       `json:"mirrorOf"`
	AvailableModes         []string     `json:"availableModes"`
	ColorManagementPreset  string       `json:"colorManagementPreset"`
	SDRBrightness          float64      `json:"sdrBrightness"`
	SDRSaturation          float64      `json:"sdrSaturation"`
	SDRMinLuminance        float64      `json:"sdrMinLuminance"`
	SDRMaxLuminance        float64      `json:"sdrMaxLuminance"`
	HardwareCursorsInUse   bool         `json:"hardwareCursorsInUse"`
}

// Monitors wraps hyprctl -j monitors: every output Hyprland knows.
func Monitors(ctx context.Context, opts ...Option) ([]Monitor, error) {
	return query[[]Monitor](ctx, "j/monitors", opts)
}
