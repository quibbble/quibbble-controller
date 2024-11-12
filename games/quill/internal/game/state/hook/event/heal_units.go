package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const (
	HealUnitsEvent = "HealUnits"
)

type HealUnitsArgs struct {
	Amount      int
	ChooseUnits parse.Choose
}

func HealUnitsAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*HealUnitsArgs)
	choose, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), a.ChooseUnits.Type, a.ChooseUnits.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, choice := range choices {
		if choice.Type() != en.UnitUUID {
			return en.ErrInvalidUUIDType(choice, en.UnitUUID)
		}
		event, err := NewEvent(state.Gen.New(en.EventUUID), HealUnitEvent, HealUnitArgs{
			Amount: a.Amount,
			ChooseUnit: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: choice,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
