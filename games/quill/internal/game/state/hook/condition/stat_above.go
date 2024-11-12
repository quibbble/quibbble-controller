package condition

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

const StatAboveCondition = "StatAbove"

type StatAboveArgs struct {
	Stat       string
	Amount     int
	ChooseCard parse.Choose
}

func PassStatAbove(c *Condition, ctx context.Context, engine *en.Engine, state *st.State) (bool, error) {
	p := c.GetArgs().(*StatAboveArgs)

	choice, err := ch.GetChoice(ctx, p.ChooseCard, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}

	card := state.GetCard(choice)
	if card == nil {
		return false, errors.ErrNilInterface
	}

	if p.Stat == cd.CostStat {
		return p.Amount < card.GetCost(), nil
	} else if choice.Type() == en.UnitUUID {
		unit := card.(*cd.UnitCard)
		switch p.Stat {
		case cd.AttackStat:
			return p.Amount < maths.MaxInt(unit.Attack, 0), nil
		case cd.HealthStat:
			return p.Amount < maths.MaxInt(unit.Health, 0), nil
		case cd.CooldownStat:
			return p.Amount < maths.MaxInt(unit.Cooldown, 0), nil
		case cd.BaseCooldownStat:
			return p.Amount < maths.MaxInt(unit.BaseCooldown, 0), nil
		case cd.MovementStat:
			return p.Amount < maths.MaxInt(unit.Movement, 0), nil
		case cd.BaseMovementStat:
			return p.Amount < maths.MaxInt(unit.BaseMovement, 0), nil
		}
	}
	return false, errors.Errorf("'%s' is not a valid stat", p.Stat)
}
