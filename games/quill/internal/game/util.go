package game

import (
	"math/rand"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	hk "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	cn "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/condition"
	ev "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/event"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

func NewDummyCard(id string) (st.ICard, error) {
	gen := uuid.NewGen(rand.New(rand.NewSource(0)))
	engineBuilders := en.Builders{
		BuildCondition: cn.NewCondition,
		BuildEvent:     ev.NewEvent,
		BuildHook:      hk.NewHook,
		BuildChoose:    ch.NewChoose,
	}
	cardBuilders := cd.Builders{
		Builders:   engineBuilders,
		BuildTrait: tr.NewTrait,
		Gen:        gen,
	}
	return cd.NewCard(&cardBuilders, id, uuid.Nil, false)
}
