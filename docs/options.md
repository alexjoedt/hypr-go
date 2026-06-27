# Options

Every function that touches a socket accepts trailing `Option` values. They
travel with the call and are resolved on each one, so there is nothing to
configure up front.

| Option               | Applies to               | Default                                              |
| -------------------- | ------------------------ | ---------------------------------------------------- |
| `WithSocketDir`      | all                      | `$XDG_RUNTIME_DIR/hypr/$HYPRLAND_INSTANCE_SIGNATURE` |
| `WithDialer`         | all                      | zero `net.Dialer`                                    |
| `WithDialTimeout`    | all                      | 5s                                                   |
| `WithCommandTimeout` | Dispatch, Request, Batch | 3s                                                   |
| `WithReadTimeout`    | Listen                   | 15s                                                  |
| `WithLogger`         | all                      | `slog.Default()`                                     |

```go
hypr.Dispatch(ctx, cmd, hypr.WithCommandTimeout(time.Second))

l, err := hypr.Listen(ctx,
    hypr.WithSocketDir("/run/user/1000/hypr/abc"),
    hypr.WithLogger(logger),
)
```

`WithDialer` takes anything with `DialContext(ctx, network, address)`.
`*net.Dialer` satisfies it. Tests use it to point the package at a fake
socket.

`WithReadTimeout` bounds a single read from the event socket. A read that
times out is not an error. The loop sets a fresh deadline and reads again,
so the value mostly controls how fast a cancelled context is noticed.
