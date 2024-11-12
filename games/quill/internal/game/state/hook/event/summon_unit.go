package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const SummonUnitEvent = "SummonUnit"

type SummonUnitArgs struct {
	ChoosePlayer parse.Choose
	ChooseID     parse.Choose
	ChooseTile   parse.Choose
	InPlayRange  bool
}

func SummonUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*SummonUnitArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	idChoice, err := ch.GetChoice(ctx, a.ChooseID, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	tileChoice, err := ch.GetTileChoice(ctx, a.ChooseTile, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	unit, err := state.BuildCard(string(idChoice), playerChoice, true)
	if err != nil {
		return errors.Wrap(err)
	}
	tX, tY, err := state.Board.GetTileXY(tileChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	if state.Board.XYs[tX][tY].Unit != nil {
		return errors.Errorf("unit '%s' cannot be placed on a full tile", unit.GetUUID())
	}
	if a.InPlayRange {
		min, max := state.Board.GetPlayableRowRange(playerChoice)
		if tY < min || tY > max {
			return errors.Errorf("unit '%s' must be placed within rows %d to %d", unit.GetUUID(), min, max)
		}
	}

	// haste trait check
	if len(unit.GetTraits(tr.HasteTrait)) > 0 {
		unit.(*cd.UnitCard).Cooldown = 0
	}

	// surge trait check
	for range unit.GetTraits(tr.SurgeTrait) {
		unit.(*cd.UnitCard).Attack += state.Mana[playerChoice].Amount
	}

	state.Board.XYs[tX][tY].Unit = unit

	// friends/enemies trait check
	if err := friendsTraitCheck(e, engine, state); err != nil {
		return errors.Wrap(err)
	}
	if err := enemiesTraitCheck(e, engine, state); err != nil {
		return errors.Wrap(err)
	}

	return nil
}
