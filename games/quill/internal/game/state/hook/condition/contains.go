package condition

import (
	"context"
	"slices"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const ContainsCondition = "Contains"

type ContainsArgs struct {
	ChooseChain parse.Choose
	Choose      parse.Choose
}

func PassContains(c *Condition, ctx context.Context, engine *en.Engine, state *st.State) (bool, error) {
	p := c.GetArgs().(*ContainsArgs)
	choices, err := ch.GetChoices(ctx, p.ChooseChain, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	choice, err := ch.GetChoice(ctx, p.Choose, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	return slices.Contains(choices, choice), nil
}
