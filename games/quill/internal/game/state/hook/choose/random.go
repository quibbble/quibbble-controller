package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const RandomChoice = "Random"

type RandomArgs struct {
	Choose parse.Choose
}

func RetrieveRandom(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*RandomArgs)
	choose, err := NewChoose(state.Gen.New(en.ChooseUUID), r.Choose.Type, r.Choose.Args)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(choices) == 0 {
		return choices, nil
	}
	return []uuid.UUID{choices[state.Rand.Intn(len(choices))]}, nil
}
