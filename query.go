package hypr

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// query sends a JSON request ("j/clients") through Request and unmarshals the
// reply into T. A reply that is not JSON, such as "unknown request", comes
// back as an error carrying its first line.
func query[T any](ctx context.Context, request string, opts []Option) (T, error) {
	var v T

	reply, err := Request(ctx, request, opts...)
	if err != nil {
		return v, err
	}

	if err := json.Unmarshal(reply, &v); err != nil {
		line, _, _ := strings.Cut(strings.TrimSpace(string(reply)), "\n")
		return v, fmt.Errorf("%s: %w: %q", request, err, line)
	}

	return v, nil
}

// WorkspaceRef is the short workspace reference embedded in clients and
// monitors.
type WorkspaceRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
