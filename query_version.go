package hypr

import "context"

// VersionInfo is what hyprctl -j version reports about the running Hyprland
// and the libraries it was built against. It is not called Version because
// the function that fetches it is.
type VersionInfo struct {
	Branch             string   `json:"branch"`
	Commit             string   `json:"commit"`
	Version            string   `json:"version"`
	Dirty              bool     `json:"dirty"`
	CommitMessage      string   `json:"commit_message"`
	CommitDate         string   `json:"commit_date"`
	Tag                string   `json:"tag"`
	Commits            string   `json:"commits"`
	BuildAquamarine    string   `json:"buildAquamarine"`
	BuildHyprlang      string   `json:"buildHyprlang"`
	BuildHyprutils     string   `json:"buildHyprutils"`
	BuildHyprcursor    string   `json:"buildHyprcursor"`
	BuildHyprgraphics  string   `json:"buildHyprgraphics"`
	SystemAquamarine   string   `json:"systemAquamarine"`
	SystemHyprlang     string   `json:"systemHyprlang"`
	SystemHyprutils    string   `json:"systemHyprutils"`
	SystemHyprcursor   string   `json:"systemHyprcursor"`
	SystemHyprgraphics string   `json:"systemHyprgraphics"`
	ABIHash            string   `json:"abiHash"`
	Flags              []string `json:"flags"`
}

// Version wraps hyprctl -j version.
func Version(ctx context.Context, opts ...Option) (*VersionInfo, error) {
	v, err := query[VersionInfo](ctx, "j/version", opts)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
