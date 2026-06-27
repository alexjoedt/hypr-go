# Events

The event socket is a stream, so it gets one stateful type. `Listen` dials
it and returns a `*Listener` that owns the connection until `Close`.

```go
l, err := hypr.Listen(ctx)
if err != nil {
    return err
}
defer l.Close()
```

The context given to `Listen` bounds the dial only. The read loop takes its
own context in `Run` or `Events`.

## Typed handlers

`On` registers a handler for exactly one event type and returns a function
that removes it again:

```go
off := hypr.On(l, func(e *hypr.ActiveWindowV2Event) {
    fmt.Println("focus:", e.Address)
})
defer off()
```

Handlers run inside the read loop, in registration order, before the event
is yielded to the iterator.

## Driving the loop

Pick one of two ways.

`Run` drives the loop and fires handlers until the context ends, `Close` is
called or the socket closes:

```go
err := l.Run(ctx)
```

`Events` yields every event after the handlers ran for it:

```go
for e := range l.Events(ctx) {
    switch e := e.(type) {
    case *hypr.WorkspaceV2Event:
        fmt.Println(e.WorkspaceID, e.WorkspaceName)
    case *hypr.UnknownEvent:
        fmt.Println("unknown:", e.Name(), e.Data())
    }
}
```

`Err` says why the last loop stopped: nil on a clean end or when you broke
out of the range, the context error on cancel, `ErrClosed` after `Close`,
else the read error.

## Unknown events

An event this version has no type for, or a payload that does not parse, is
an `*UnknownEvent`. It is a value in the stream, not a dropped line. `Name`
returns the part before `>>`, `Data` the part after.

## Parsing offline

`Parse(line)` turns one event line into an `Event` without a socket. It is
useful for tests and for replaying recorded streams.

## Timeouts

Hyprland can stay silent for a long time. A read that times out is not an
error: the loop sets a fresh deadline and reads again. See `WithReadTimeout`
in [options](options.md).
