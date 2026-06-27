# Contributing

## Layout

hypr-go is one package. Do not add subpackages for commands, queries or
events: the three families share `WindowAddress`, `WorkspaceRef`, the
selectors and `Option`, and every real program uses all of them. The flat
surface is kept navigable by naming, not by packages.

## Naming

The caller-facing rule is in [doc.go](doc.go) under "Finding a name". The
tiebreaks contributors need:

- **Dispatcher.** Go name is the Lua path in CamelCase, nothing else.
  Top-level dispatchers stay bare (`Focus`, `Pass`, `DPMS`), no family prefix.
  An overloaded dispatcher gets one struct per mode with the mode appended
  (`FocusWindow`, `FocusMonitor`). Table dispatchers are structs whose fields
  mirror the Lua keys and whose zero value is Hyprland's default. Positional
  and no-arg dispatchers are functions returning `Command`. No per-dispatcher
  functional options.
- **Event.** Go type is the event name plus `Event`, case adjusted
  (`workspacev2` is `WorkspaceV2Event`). Every event has a constant, a
  registry entry and a row in `TestParseEvent`.
- **Query.** Function is the hyprctl name. A list reply returns the singular
  noun (`Clients` returns `[]Client`). An object reply returns the name plus
  `Info` (`Version` returns `*VersionInfo`), because a function and a type
  cannot share a name.
- **Selector.** Typed string per noun, constructors `NounBy…`
  (`WorkspaceByID`, not `WorkspaceID`), constants `Noun…` for the fixed
  values (`WorkspaceNext`). `NounRelative` and `MonitorInDirection` are the
  two constructors without a key.
- **Enum.** Constants start with the type name (`ActionToggle`,
  `DirectionLeft`) or the type name minus its last word when that reads
  better (`KeyDown` for `KeyState`, `ScreenCastMonitor` for
  `ScreenCastOwner`). Never a bare `Toggle` or `Left`.
- **`EmitEvent`.** The only dispatcher not named after its Lua path,
  `hl.dsp.event`, because `Event` is the interface. Do not add a second
  exception; prefix by family instead.
- **`Client` is not `Window`.** The data struct follows hyprctl, the rest
  follows the dispatcher and event names. Do not alias.

`TestDispatcherCoverage`, `TestParseEvent` and `TestEnumConstantNames` check
the dispatcher, event and enum shapes. Queries are few enough to check by eye.

## Adding a dispatcher

1. Add it to `testdata/dispatchers.txt` if Hyprland documents it.
2. Add the struct or function to the family file (`dispatch_window.go`, …).
3. Add a row to the family's wire-string table in the test file. Coverage is
   derived from those tables.

## Adding an event

1. Add the constant to `event.go` and the parser to the registry.
2. Add the type to `event_types.go`.
3. Add rows to `TestParseEvent` and a line to `testdata/event_stream.txt`.

## Adding a query

1. Add the result type and function to a `query_<noun>.go` file.
2. Add a fixture under `testdata/query/` and a test.

## Checks

```sh
make audit   # vet, lint, race tests, module verify
```

Commits follow Conventional Commits. Do not add `Co-Authored-By` trailers.
