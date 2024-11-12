package state

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type BuildCard func(id string, player uuid.UUID, token bool) (ICard, error)

type ICard interface {
	GetID() string
	GetUUID() uuid.UUID
	GetPlayer() uuid.UUID
	GetCost() int
	SetCost(cost int)
	GetInit() parse.ICard
	GetEvents() []en.IEvent
	GetHooks() []en.IHook
	Playable(engine en.IEngine, state en.IState) (bool, error)
	GetTargets() []en.IChoose
	NextTargets(ctx context.Context, engine en.IEngine, state en.IState) ([]uuid.UUID, error)
	GetTraits(typ string) []ITrait
	AddTrait(trait ITrait) error
	RemoveTrait(trait uuid.UUID) error
}

type BuildTrait func(uuid uuid.UUID, createdBy *uuid.UUID, typ string, args interface{}) (ITrait, error)

type ITrait interface {
	GetUUID() uuid.UUID
	GetType() string
	GetArgs() interface{}
	GetCreatedBy() *uuid.UUID
	Add(card ICard) error
	Remove(card ICard) error
}
