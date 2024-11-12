package condition

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const MatchCondition = "Match"

type MatchArgs struct {
	ChooseA parse.Choose
	ChooseB parse.Choose
}

func PassMatch(c *Condition, ctx context.Context, engine *en.Engine, state *st.State) (bool, error) {
	p := c.GetArgs().(*MatchArgs)
	a, err := choose.GetChoice(ctx, p.ChooseA, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	b, err := choose.GetChoice(ctx, p.ChooseB, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	return a == b, nil
}
