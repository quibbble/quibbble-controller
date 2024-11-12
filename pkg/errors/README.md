# Errors

This package creates a custom error for reporting stacktraces useful for debugging complex game logic. 

```
go-quill/cards/spells/S000_test.go:15: failed to build event from 'DamageUnitFromTarget' raw event
    spells.play
        go-xalpha/cards/spells/util.go:42
    event_system.(*EventSystem).Do
        go-xalpha/internal/event_system/event_system.go:43
    event_system.(*EventSystem).do
        go-xalpha/internal/event_system/event_system.go:70
    event_builders.BuildEvents
        go-xalpha/internal/event_system/events/event_builders/1_event_builder.go:58
```

Create a new error:
```go
import module github.com/quibbble/quibbble-controller/pkg/errors

err := errors.Errorf("this is an %s", "example")
```

Wrap an existing error:
```go
import module github.com/quibbble/quibbble-controller/pkg/errors

if err := example(); err != nil {
    return errors.Wrap(err)
}
```
