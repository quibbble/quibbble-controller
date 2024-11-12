package choose

import (
	"context"
	"slices"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/maths"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const RangedChoice = "Ranged"

type RangedArgs struct {
	Types            []string
	Range            int
	ChooseUnitOrTile parse.Choose
}

func RetrieveRanged(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*RangedArgs)
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

	unitX, unitY, err := state.Board.GetUnitXY(choice)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	ranged := make([]uuid.UUID, 0)
	for x, col := range state.Board.XYs {
		for y, tile := range col {
			if maths.AbsInt(unitX-x) <= r.Range && maths.AbsInt(unitY-y) <= r.Range {
				if slices.Contains(r.Types, "Tile") {
					ranged = append(ranged, tile.UUID)
				} else if tile.Unit != nil {
					unit := tile.Unit.(*cd.UnitCard)
					if slices.Contains(r.Types, cd.Unit) || slices.Contains(r.Types, unit.Type) {
						ranged = append(ranged, unit.UUID)
					}
				}
			}
		}
	}
	return ranged, nil
}
