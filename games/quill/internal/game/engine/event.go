package engine

import (
	"context"

	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type BuildEvent func(uuid uuid.UUID, typ string, args interface{}) (IEvent, error)

// Events make a change to the state
type IEvent interface {
	GetUUID() uuid.UUID
	GetType() string
	GetArgs() interface{}
	Affect(ctx context.Context, engine IEngine, state IState) error
}
