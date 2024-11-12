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
	ModifyUnitEvent = "ModifyUnit"
)

type ModifyUnitArgs struct {
	ChooseUnit parse.Choose
	Stat       string
	Amount     int
}

func ModifyUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*ModifyUnitArgs)
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)

	switch a.Stat {
	case cd.AttackStat:
		unit.Attack = unit.Attack + a.Amount
	case cd.HealthStat:
		unit.Health = unit.Health + a.Amount
	case cd.MovementStat:
		unit.Movement = maths.MaxInt(0, unit.Movement+a.Amount)
	case cd.CooldownStat:
		unit.Cooldown = maths.MaxInt(0, unit.Cooldown+a.Amount)
	default:
		return errors.Errorf("'%s' is not a stat that may be modified", a.Stat)
	}

	return nil
}
