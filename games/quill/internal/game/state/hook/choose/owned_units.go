package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const OwnedUnitsChoice = "OwnedUnits"

type OwnedUnitsArgs struct {
	ChoosePlayer parse.Choose
}

func RetrieveOwnedUnits(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*OwnedUnitsArgs)
	choose, err := NewChoose(state.Gen.New(en.ChooseUUID), r.ChoosePlayer.Type, r.ChoosePlayer.Args)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(choices) != 1 || choices[0].Type() != en.PlayerUUID {
		return nil, errors.ErrInvalidSliceLength
	}
	owned := make([]uuid.UUID, 0)
	for _, tile := range state.Board.UUIDs {
		if tile.Unit != nil && tile.Unit.GetPlayer() == choices[0] {
			owned = append(owned, tile.Unit.GetUUID())
		}
	}
	return owned, nil
}
