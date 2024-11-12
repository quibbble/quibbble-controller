package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/maths"
)

const (
	RefreshMovementEvent = "RefreshMovement"
)

type RefreshMovementArgs struct {
	ChooseUnits parse.Choose
}

func RefreshMovementAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*RefreshMovementArgs)
	choose, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), a.ChooseUnits.Type, a.ChooseUnits.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, choice := range choices {
		x, y, err := state.Board.GetUnitXY(choice)
		if err != nil {
			return errors.Wrap(err)
		}
		unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)

		event, err := NewEvent(state.Gen.New(en.EventUUID), ModifyUnitEvent, ModifyUnitArgs{
			ChooseUnit: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: unit.GetUUID(),
				},
			},
			Stat:   cd.MovementStat,
			Amount: maths.MaxInt(0, unit.BaseMovement) - unit.Movement,
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
