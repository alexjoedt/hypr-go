// Command example shows the ways to talk to Hyprland: Dispatch for a one-off
// command, ActiveWindow for a typed query, On for typed handlers driven by
// Run, and, with -events, the Events iterator instead of Run.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	hypr "github.com/alexjoedt/hypr-go"
)

func main() {
	iterate := flag.Bool("events", false, "range over Events instead of Run with handlers")
	flag.Parse()

	if err := run(*iterate); err != nil {
		fmt.Fprintf(os.Stderr, "exiting with an error: %v\n", err)
		os.Exit(1)
	}
}

func run(iterate bool) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// A script: one call, nothing to construct, nothing to close. Every
	// dispatcher has a typed constructor; RawCommand is the escape hatch for
	// ones a newer Hyprland adds.
	if err := hypr.Dispatch(ctx, hypr.NoOp()); err != nil {
		var re *hypr.ReplyError
		if !errors.As(err, &re) {
			return err
		}
		log.Printf("hyprland refused %q: %v", re.Request, re.Errors)
	}

	// A query: the reply of hyprctl -j as a typed struct.
	if w, err := hypr.ActiveWindow(ctx); err != nil {
		return err
	} else if w == nil {
		fmt.Println("no focused window")
	} else {
		fmt.Printf("focused: %s (%s)\n", w.Title, w.Class)
	}

	// A daemon: one Listener owns the event socket until Close.
	l, err := hypr.Listen(ctx)
	if err != nil {
		return err
	}
	defer l.Close()

	// Typed handlers see exactly one event type each.
	unsubscribe := hypr.On(l, func(e *hypr.ActiveWindowV2Event) {
		fmt.Println("focus:", e.Address)
	})
	defer unsubscribe()

	if !iterate {
		// Run drives the loop and fires the handlers until ctx ends.
		if err := l.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}
		return nil
	}

	// The iterator sees everything, after the handlers ran for it. Unknown
	// events are values, not dropped lines.
	for e := range l.Events(ctx) {
		switch e := e.(type) {
		case *hypr.WorkspaceV2Event:
			fmt.Printf("workspace: %d %s\n", e.WorkspaceID, e.WorkspaceName)
		case *hypr.UnknownEvent:
			fmt.Printf("unknown event %q: %s\n", e.Name(), e.Data())
		default:
			fmt.Println(e.Name())
		}
	}

	if err := l.Err(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
