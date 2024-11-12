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

const CodexChoice = "Codex"

var codexXYs = [][]int{{0, 1}, {0, -1}, {-1, 0}, {1, 0}, {-1, 1}, {1, -1}, {-1, -1}, {1, 1}}

type CodexArgs struct {
	Types            []string
	Codex            string
	ChooseUnitOrTile parse.Choose
}

func RetrieveCodex(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*CodexArgs)
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
	codex := make([]uuid.UUID, 0)
	for i := 0; i < len(r.Codex); i++ {
		if r.Codex[i] == '1' {
			x, y := x+codexXYs[i][0], y+codexXYs[i][1]
			if x < 0 || x >= st.Cols || y < 0 || y >= st.Rows {
				continue
			}
			tile := state.Board.XYs[x][y]
			if slices.Contains(r.Types, "Tile") {
				codex = append(codex, tile.UUID)
			} else if tile.Unit != nil {
				unit := tile.Unit.(*cd.UnitCard)
				if slices.Contains(r.Types, cd.Unit) || slices.Contains(r.Types, unit.Type) {
					codex = append(codex, unit.UUID)
				}
			}
		}
	}
	return codex, nil
}
