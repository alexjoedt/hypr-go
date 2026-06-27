# hypr-go

[![ci](https://github.com/alexjoedt/hypr-go/actions/workflows/ci.yml/badge.svg)](https://github.com/alexjoedt/hypr-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexjoedt/hypr-go.svg)](https://pkg.go.dev/github.com/alexjoedt/hypr-go)

Go library for the Hyprland IPC sockets. Typed dispatchers and `hyprctl -j`
queries over the command socket, typed events from the event socket.
Standard library only. Needs Hyprland 0.56.2 or newer.

## Install

```sh
go get github.com/alexjoedt/hypr-go
```

## Usage

```go
import hypr "github.com/alexjoedt/hypr-go"

// Send a dispatcher.
err := hypr.Dispatch(ctx, hypr.FocusWorkspace{Workspace: hypr.WorkspaceByID(3)})

// Run a query and get a typed struct back.
w, err := hypr.ActiveWindow(ctx)

// Listen for events.
l, err := hypr.Listen(ctx)
defer l.Close()

hypr.On(l, func(e *hypr.ActiveWindowV2Event) { fmt.Println("focus:", e.Address) })

for e := range l.Events(ctx) {
    fmt.Println(e.Name())
}
```

Every dispatcher of Hyprland 0.56.2 has a typed constructor. Selectors are
built with functions such as `hypr.WindowByClass("^kitty$")` or
`hypr.MonitorByName("DP-1")`. A reply Hyprland rejects comes back as a
`*hypr.ReplyError`. An event this version has no type for comes back as an
`*hypr.UnknownEvent`.

For anything newer than this library knows, use `hypr.RawCommand` for
dispatchers and `hypr.Request` for other hyprctl requests.

`examples/main.go` is a runnable version of the snippet above:

```sh
go run ./examples
go run ./examples -events
```

## Documentation

Short guides are in [docs/](docs/index.md). The API reference is on
[pkg.go.dev](https://pkg.go.dev/github.com/alexjoedt/hypr-go).

## Development

```sh
make lint    # golangci-lint
make test    # race tests with coverage
make audit   # vet, lint, race tests, module verify
```

## License

MIT, see [LICENSE](LICENSE).
