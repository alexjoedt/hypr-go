package hypr_test

import (
	"context"
	"fmt"
	"log"

	hypr "github.com/alexjoedt/hypr-go"
)

// A dispatcher is a value; Dispatch dials, sends it and closes.
func ExampleDispatch() {
	ctx := context.Background()

	cmd := hypr.WindowMoveToWorkspace{Workspace: hypr.WorkspaceByID(2), Window: hypr.WindowByClass("^kitty$")}
	if err := hypr.Dispatch(ctx, cmd); err != nil {
		log.Fatal(err)
	}
}

// A query returns the reply of hyprctl -j as a typed value.
func ExampleActiveWindow() {
	ctx := context.Background()

	w, err := hypr.ActiveWindow(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if w == nil {
		fmt.Println("nothing focused")
		return
	}
	fmt.Println(w.Class, w.Title)
}

// On registers a handler for one event type; Run drives the loop.
func ExampleOn() {
	ctx := context.Background()

	l, err := hypr.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	unsubscribe := hypr.On(l, func(e *hypr.WorkspaceV2Event) {
		fmt.Println("workspace", e.WorkspaceID, e.WorkspaceName)
	})
	defer unsubscribe()

	if err := l.Run(ctx); err != nil {
		fmt.Println("listener stopped:", err)
	}
}
