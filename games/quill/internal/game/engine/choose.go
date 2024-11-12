package engine

import (
	"context"

	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type BuildChoose func(uuid uuid.UUID, typ string, args interface{}) (IChoose, error)

type IChoose interface {
	Retrieve(ctx context.Context, engine IEngine, state IState) ([]uuid.UUID, error)
}
