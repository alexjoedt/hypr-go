package hypr

import "context"

// Client is one window as hyprctl -j clients reports it. Hyprland says client
// here and window everywhere else; [Clients] and [ActiveWindow] return it, the
// [WindowSelector] and Window dispatchers address it.
type Client struct {
	Address               WindowAddress   `json:"address"`
	Mapped                bool            `json:"mapped"`
	Hidden                bool            `json:"hidden"`
	Visible               bool            `json:"visible"`
	AcceptsInput          bool            `json:"acceptsInput"`
	At                    [2]int          `json:"at"`
	Size                  [2]int          `json:"size"`
	Workspace             WorkspaceRef    `json:"workspace"`
	Floating              bool            `json:"floating"`
	Monitor               int             `json:"monitor"`
	Class                 string          `json:"class"`
	Title                 string          `json:"title"`
	InitialClass          string          `json:"initialClass"`
	InitialTitle          string          `json:"initialTitle"`
	PID                   int             `json:"pid"`
	XWayland              bool            `json:"xwayland"`
	Pinned                bool            `json:"pinned"`
	PinFullscreened       bool            `json:"pinFullscreened"`
	Fullscreen            FullscreenState `json:"fullscreen"`
	FullscreenClient      FullscreenState `json:"fullscreenClient"`
	FullscreenHandler     string          `json:"fullscreenHandler"`
	AllowedOverFullscreen bool            `json:"allowedOverFullscreen"`
	Grouped               []WindowAddress `json:"grouped"`
	Tags                  []string        `json:"tags"`
	Swallowing            WindowAddress   `json:"swallowing"`
	FocusHistoryID        int             `json:"focusHistoryID"`
	InhibitingIdle        bool            `json:"inhibitingIdle"`
	XDGTag                string          `json:"xdgTag"`
	XDGDescription        string          `json:"xdgDescription"`
	ContentType           string          `json:"contentType"`
	TearingHint           bool            `json:"tearingHint"`
	StableID              string          `json:"stableId"`
}

// Clients wraps hyprctl -j clients: every window Hyprland manages, one
// [Client] each. See [ActiveWindow] for the focused one.
func Clients(ctx context.Context, opts ...Option) ([]Client, error) {
	return query[[]Client](ctx, "j/clients", opts)
}

// ActiveWindow wraps hyprctl -j activewindow: the focused window as a
// [Client], or nil when nothing is focused. See [Clients] for all of them.
func ActiveWindow(ctx context.Context, opts ...Option) (*Client, error) {
	// Hyprland replies {} when nothing is focused, which leaves the address
	// empty.
	c, err := query[Client](ctx, "j/activewindow", opts)
	if err != nil || c.Address == "" {
		return nil, err
	}

	return &c, nil
}
