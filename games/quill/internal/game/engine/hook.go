package engine

import (
	"context"

	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type When string

const (
	Before When = "Before"
	After  When = "After"
)

type BuildHook func(uuid, cardUUID uuid.UUID, when string, types []string, conditions []ICondition, events []IEvent, reuse []ICondition) (IHook, error)

// Hooks are always registered by a card
type IHook interface {
	GetUUID() uuid.UUID
	GetCardUUID() uuid.UUID
	GetTypes() []string
	Trigger(when When, typ string) bool
	Pass(ctx context.Context, engine IEngine, state IState) (bool, error)
	Events() []IEvent
	Reuse(ctx context.Context, engine IEngine, state IState) (bool, error)
}
