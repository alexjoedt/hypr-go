package hypr

import (
	"context"
	"reflect"
	"testing"
)

func TestWorkspaces(t *testing.T) {
	t.Parallel()

	workspaces, err := Workspaces(context.Background(), fixtureOpts(t, "workspaces.json")...)
	if err != nil {
		t.Fatalf("Workspaces: %v", err)
	}
	if len(workspaces) != 7 {
		t.Fatalf("got %d workspaces, want 7", len(workspaces))
	}

	special := Workspace{
		ID: -98, Name: "special:magic", Monitor: "eDP-1", MonitorID: 0, Windows: 1,
		LastWindow: "559d0cde5e60", LastWindowTitle: "Window", TiledLayout: "dwindle",
	}
	if !reflect.DeepEqual(workspaces[0], special) {
		t.Errorf("workspaces[0]\n got %+v\nwant %+v", workspaces[0], special)
	}

	numbered := Workspace{
		ID: 2, Name: "2", Monitor: "eDP-1", MonitorID: 0, Windows: 1,
		LastWindow: "559d0cd98040", LastWindowTitle: "Window", TiledLayout: "dwindle",
	}
	if !reflect.DeepEqual(workspaces[1], numbered) {
		t.Errorf("workspaces[1]\n got %+v\nwant %+v", workspaces[1], numbered)
	}
}

func TestActiveWorkspace(t *testing.T) {
	t.Parallel()

	w, err := ActiveWorkspace(context.Background(), fixtureOpts(t, "activeworkspace.json")...)
	if err != nil {
		t.Fatalf("ActiveWorkspace: %v", err)
	}
	if w == nil || w.ID != 3 || w.Name != "3" || w.Windows != 2 || w.LastWindow != "559d0ceb9260" {
		t.Errorf("ActiveWorkspace = %+v", w)
	}
}
