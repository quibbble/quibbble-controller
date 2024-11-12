package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const ModifyUnitsEvent = "ModifyUnits"

type ModifyUnitsArgs struct {
	ChooseUnits parse.Choose
	Stat        string
	Amount      int
}

func ModifyUnitsAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*ModifyUnitsArgs)
	unitChoices, err := ch.GetChoices(ctx, a.ChooseUnits, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, choice := range unitChoices {
		event, err := NewEvent(state.Gen.New(en.EventUUID), ModifyUnitEvent, ModifyUnitArgs{
			ChooseUnit: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: choice,
				},
			},
			Stat:   a.Stat,
			Amount: a.Amount,
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(ctx, event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
