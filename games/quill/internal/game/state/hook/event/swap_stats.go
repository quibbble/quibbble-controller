package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const SwapStatsEvent = "SwapStats"

type SwapStatsArgs struct {
	Stat        string
	ChooseCardA parse.Choose
	ChooseCardB parse.Choose
}

func SwapStatsAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*SwapStatsArgs)
	choiceA, err := ch.GetChoice(ctx, a.ChooseCardA, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	choiceB, err := ch.GetChoice(ctx, a.ChooseCardB, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	cardA := state.GetCard(choiceA)
	if cardA == nil {
		return errors.ErrNilInterface
	}
	cardB := state.GetCard(choiceB)
	if cardB == nil {
		return errors.ErrNilInterface
	}

	if a.Stat == cd.CostStat {
		costA := cardA.GetCost()
		costB := cardB.GetCost()
		cardA.SetCost(costB)
		cardB.SetCost(costA)
		return nil
	} else if cardA.GetUUID().Type() == en.UnitUUID && cardB.GetUUID().Type() == en.UnitUUID {
		unitA := cardA.(*cd.UnitCard)
		unitB := cardB.(*cd.UnitCard)
		switch a.Stat {
		case cd.AttackStat:
			attackA := unitA.Attack
			attackB := unitB.Attack
			unitA.Attack = attackB
			unitB.Attack = attackA
		case cd.HealthStat:
			healthA := unitA.Health
			healthB := unitB.Health
			unitA.Health = healthB
			unitB.Health = healthA
		case cd.MovementStat:
			moveA := unitA.Movement
			moveB := unitB.Movement
			unitA.Movement = moveB
			unitB.Movement = moveA
		case cd.CooldownStat:
			coolA := unitA.Cooldown
			coolB := unitB.Cooldown
			unitA.Cooldown = coolB
			unitB.Cooldown = coolA
		}
		return nil
	}
	return errors.Errorf("'%s' cannot be swapped between '%s' and '%s'", a.Stat, choiceA, choiceB)
}
