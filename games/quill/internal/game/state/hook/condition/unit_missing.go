package condition

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const UnitMissingCondition = "UnitMissing"

type UnitMissingArgs struct {
	ChooseUnit parse.Choose
}

func PassUnitMissing(c *Condition, ctx context.Context, engine *en.Engine, state *st.State) (bool, error) {
	p := c.GetArgs().(*UnitMissingArgs)

	unitChoice, err := ch.GetUnitChoice(ctx, p.ChooseUnit, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}

	_, _, err = state.Board.GetUnitXY(unitChoice)
	return err != nil, nil
}
