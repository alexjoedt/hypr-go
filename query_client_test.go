package hypr

import (
	"context"
	"reflect"
	"testing"
)

func TestClients(t *testing.T) {
	t.Parallel()

	opts := fixtureOpts(t, "clients.json")

	clients, err := Clients(context.Background(), opts...)
	if err != nil {
		t.Fatalf("Clients: %v", err)
	}
	if len(clients) != 7 {
		t.Fatalf("got %d clients, want 7", len(clients))
	}

	// The grouped, tagged editor on workspace 4.
	want := Client{
		Address: "559d0ced3cb0", Mapped: true, Visible: true, AcceptsInput: true,
		At: [2]int{2232, 61}, Size: [2]int{1856, 1358},
		Workspace: WorkspaceRef{ID: 4, Name: "4"},
		Monitor:   1, Class: "codium", Title: "Editor", InitialClass: "codium", InitialTitle: "Editor",
		PID: 739856, FullscreenHandler: "default", AllowedOverFullscreen: true,
		Grouped: []WindowAddress{"559d0ced3cb0", "559d0ced4f10"}, Tags: []string{"urgent"},
		Swallowing: "0", FocusHistoryID: 4, ContentType: "none", StableID: "18000065",
	}
	if !reflect.DeepEqual(clients[0], want) {
		t.Errorf("clients[0]\n got %+v\nwant %+v", clients[0], want)
	}

	// The notes app on the special workspace.
	if c := clients[2]; c.Class != "md.obsidian.Obsidian" || c.Workspace != (WorkspaceRef{ID: -98, Name: "special:magic"}) {
		t.Errorf("clients[2] = %+v, want the client on special:magic", c)
	}

	// The maximized terminal.
	if c := clients[3]; c.Fullscreen != FullscreenStateMaximized || c.FullscreenClient != FullscreenStateMaximized {
		t.Errorf("clients[3] fullscreen = %d/%d, want maximized/maximized", c.Fullscreen, c.FullscreenClient)
	}

	// The floating terminal with an empty tag.
	if c := clients[6]; !c.Floating || len(c.Tags) != 1 || c.Tags[0] != "" || c.Class != "org.wezfurlong.wezterm" {
		t.Errorf("clients[6] = %+v, want the floating terminal", c)
	}

	for i, c := range clients {
		if c.Address == "" || c.Address[:2] == "0x" {
			t.Errorf("clients[%d].Address = %q, want the hex without 0x", i, c.Address)
		}
	}
}

func TestActiveWindow(t *testing.T) {
	t.Parallel()

	t.Run("focused", func(t *testing.T) {
		t.Parallel()

		c, err := ActiveWindow(context.Background(), fixtureOpts(t, "activewindow.json")...)
		if err != nil {
			t.Fatalf("ActiveWindow: %v", err)
		}
		if c == nil {
			t.Fatal("ActiveWindow = nil, want the focused window")
		}
		if c.Address != "559d0ceb9260" || c.Class != "org.wezfurlong.wezterm" || c.Title != "Terminal" || !c.Floating {
			t.Errorf("ActiveWindow = %+v", c)
		}
		if c.Workspace != (WorkspaceRef{ID: 3, Name: "3"}) || c.At != [2]int{0, 51} {
			t.Errorf("ActiveWindow workspace/at = %+v/%v", c.Workspace, c.At)
		}
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		c, err := ActiveWindow(context.Background(), fixtureOpts(t, "activewindow_none.json")...)
		if err != nil {
			t.Fatalf("ActiveWindow: %v", err)
		}
		if c != nil {
			t.Errorf("ActiveWindow = %+v, want nil", c)
		}
	})
}
