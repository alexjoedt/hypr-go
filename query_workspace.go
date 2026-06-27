package hypr

import "context"

// Workspace is one workspace as hyprctl -j workspaces reports it. Special
// workspaces have a negative ID and a "special:" name.
type Workspace struct {
	ID              int           `json:"id"`
	Name            string        `json:"name"`
	Monitor         string        `json:"monitor"`
	MonitorID       int           `json:"monitorID"`
	Windows         int           `json:"windows"`
	HasFullscreen   bool          `json:"hasfullscreen"`
	LastWindow      WindowAddress `json:"lastwindow"`
	LastWindowTitle string        `json:"lastwindowtitle"`
	IsPersistent    bool          `json:"ispersistent"`
	TiledLayout     string        `json:"tiledLayout"`
}

// Workspaces wraps hyprctl -j workspaces: every workspace that exists.
func Workspaces(ctx context.Context, opts ...Option) ([]Workspace, error) {
	return query[[]Workspace](ctx, "j/workspaces", opts)
}

// ActiveWorkspace wraps hyprctl -j activeworkspace: the workspace of the
// focused monitor.
func ActiveWorkspace(ctx context.Context, opts ...Option) (*Workspace, error) {
	w, err := query[Workspace](ctx, "j/activeworkspace", opts)
	if err != nil {
		return nil, err
	}

	return &w, nil
}
