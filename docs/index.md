# hypr-go documentation

hypr-go is one Go package that talks to a running Hyprland over its two
unix sockets. The API follows the sockets: request and reply are functions,
the event stream is a type.

- [Commands and queries](commands.md): `Dispatch`, `Batch`, `Request` and
  the typed `hyprctl -j` queries over the command socket.
- [Events](events.md): `Listen`, `On` and `Events` over the event socket.
- [Options](options.md): timeouts, socket directory, dialer and logger.

Every name is derived from the one Hyprland uses. The rule is in the package
doc under "Finding a name"; [CONTRIBUTING.md](../CONTRIBUTING.md) has the
tiebreaks.

The full reference is on
[pkg.go.dev](https://pkg.go.dev/github.com/alexjoedt/hypr-go). A runnable
example lives in [examples/main.go](../examples/main.go).

## How the sockets are found

Every call reads `XDG_RUNTIME_DIR` and `HYPRLAND_INSTANCE_SIGNATURE` from the
environment and dials `$XDG_RUNTIME_DIR/hypr/$HYPRLAND_INSTANCE_SIGNATURE`.
If the signature is empty, the call returns `ErrNoHyprlandSocket`. Use
`WithSocketDir` to point at another directory.

## Errors

- `*ReplyError`: Hyprland answered, but not with `ok`. It carries the request
  string, the error lines and the warning lines.
- `ErrNoHyprlandSocket`: no instance signature in the environment.
- Everything else is a plain dial, read or context error.
