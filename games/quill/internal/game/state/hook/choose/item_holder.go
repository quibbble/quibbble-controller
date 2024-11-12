package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const ItemHolderChoice = "ItemHolder"

type ItemHolderArgs struct {
	ChooseItem parse.Choose
}

func RetrieveItemHolder(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*ItemHolderArgs)
	item, err := GetItemChoice(ctx, r.ChooseItem, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	holder := uuid.Nil
	for _, tile := range state.Board.UUIDs {
		if tile.Unit != nil {
			_, err := tile.Unit.(*cd.UnitCard).GetItem(item)
			if err == nil {
				holder = tile.Unit.GetUUID()
				break
			}
		}
	}
	if holder == uuid.Nil {
		return nil, errors.Errorf("'%s' not held by any unit", item)
	}
	return []uuid.UUID{holder}, nil
}
