package choose

import (
	"context"
	"slices"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const AdjacentChoice = "Adjacent"

var adjacentXYs = [][]int{{-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}}

type AdjacentArgs struct {
	Types            []string
	ChooseUnitOrTile parse.Choose
}

func RetrieveAdjacent(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*AdjacentArgs)
	choose, err := NewChoose(state.Gen.New(en.ChooseUUID), r.ChooseUnitOrTile.Type, r.ChooseUnitOrTile.Args)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(choices) != 1 {
		return nil, errors.ErrInvalidSliceLength
	}
	choice := choices[0]

	x, y, err := state.Board.GetUnitXY(choice)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	adjacent := make([]uuid.UUID, 0)
	for _, xy := range adjacentXYs {
		x, y := x+xy[0], y+xy[1]
		if x < 0 || x >= st.Cols || y < 0 || y >= st.Rows {
			continue
		}

		tile := state.Board.XYs[x][y]
		if slices.Contains(r.Types, "Tile") {
			adjacent = append(adjacent, tile.UUID)
		} else if tile.Unit != nil {
			unit := tile.Unit.(*cd.UnitCard)
			if slices.Contains(r.Types, cd.Unit) || slices.Contains(r.Types, unit.Type) {
				adjacent = append(adjacent, unit.UUID)
			}
		}
	}
	return adjacent, nil
}
