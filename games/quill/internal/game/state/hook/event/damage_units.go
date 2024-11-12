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
	DamageUnitsEvent = "DamageUnits"
)

type DamageUnitsArgs struct {
	DamageType  string
	Amount      int
	ChooseUnits parse.Choose
}

func DamageUnitsAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*DamageUnitsArgs)
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
		// it's possible the unit was killed by another affect between choosing and applying
		// if this happens then just continue
		_, _, err := state.Board.GetUnitXY(choice)
		if err != nil {
			continue
		}
		event, err := NewEvent(state.Gen.New(en.EventUUID), DamageUnitEvent, DamageUnitArgs{
			DamageType: a.DamageType,
			Amount:     a.Amount,
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
