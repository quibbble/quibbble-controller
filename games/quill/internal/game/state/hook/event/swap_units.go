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
	SwapUnitsEvent = "SwapUnits"
)

type SwapUnitsArgs struct {
	ChooseUnitA parse.Choose
	ChooseUnitB parse.Choose
}

func SwapUnitsAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*SwapUnitsArgs)
	unitAChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnitA, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	unitBChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnitB, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	aX, aY, err := state.Board.GetUnitXY(unitAChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	bX, bY, err := state.Board.GetUnitXY(unitBChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unitA := state.Board.XYs[aX][aY].Unit
	unitB := state.Board.XYs[bX][bY].Unit
	state.Board.XYs[aX][aY].Unit = unitB
	state.Board.XYs[bX][bY].Unit = unitA

	// friends/enemies trait check
	if err := friendsTraitCheck(e, engine, state); err != nil {
		return errors.Wrap(err)
	}
	if err := enemiesTraitCheck(e, engine, state); err != nil {
		return errors.Wrap(err)
	}

	return nil
}
