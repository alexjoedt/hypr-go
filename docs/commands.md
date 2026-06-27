# Commands and queries

All functions on the command socket dial, send, read the reply and close.
There is no client to construct or keep. Each call takes a context and
optional [options](options.md).

## Dispatch

`Dispatch(ctx, cmd)` runs one dispatcher. Every dispatcher of Hyprland
0.56.2 has a typed constructor.

Dispatchers with a table argument are structs. The zero value is Hyprland's
default, so you only set what you change:

```go
hypr.Dispatch(ctx, hypr.WindowClose{})
hypr.Dispatch(ctx, hypr.WindowMoveToWorkspace{
    Workspace: hypr.WorkspaceBySpecial("magic"),
    NoFollow:  true,
})
```

Dispatchers with a positional argument or no argument are functions:

```go
hypr.Dispatch(ctx, hypr.ExecCmd("kitty"))
hypr.Dispatch(ctx, hypr.NoOp())
```

Windows, workspaces and monitors are picked with selector constructors,
never with free strings:

```go
hypr.WindowByClass("^kitty$")
hypr.WindowByAddress(addr)
hypr.WorkspaceByID(3)
hypr.WorkspaceRelative(+1)
hypr.MonitorByName("DP-1")
```

A reply other than `ok` is a `*ReplyError`:

```go
if err := hypr.Dispatch(ctx, cmd); err != nil {
    var re *hypr.ReplyError
    if errors.As(err, &re) {
        log.Println(re.Request, re.Errors)
    }
}
```

## Batch

`Batch(ctx, cmds)` sends several dispatchers in one round trip. Each failed
command becomes a `*ReplyError`. Several are joined with `errors.Join`.

## RawCommand

`RawCommand` is a string that implements `Command` verbatim. Use it for a
dispatcher a newer Hyprland added:

```go
hypr.Dispatch(ctx, hypr.RawCommand("hl.dsp.new_thing()"))
```

## Queries

Queries wrap `hyprctl -j` and return typed structs:

| Function          | Returns        |
| ----------------- | -------------- |
| `Clients`         | `[]Client`     |
| `ActiveWindow`    | `*Client`      |
| `Workspaces`      | `[]Workspace`  |
| `ActiveWorkspace` | `*Workspace`   |
| `Monitors`        | `[]Monitor`    |
| `Version`         | `*VersionInfo` |

`ActiveWindow` returns nil without an error when nothing is focused.

## Request

`Request(ctx, request)` sends any hyprctl request unchanged and returns the
raw reply bytes. It covers everything that has no typed function yet:

```go
raw, err := hypr.Request(ctx, "j/devices")
```
