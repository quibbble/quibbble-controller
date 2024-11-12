package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const OwnerChoice = "Owner"

type OwnerArgs struct {
	ChooseCard parse.Choose
}

func RetrieveOwner(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*OwnerArgs)
	choice, err := GetChoice(ctx, r.ChooseCard, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	card := state.GetCard(choice)
	if card == nil {
		return nil, errors.ErrNilInterface
	}

	return []uuid.UUID{card.GetPlayer()}, nil
}
