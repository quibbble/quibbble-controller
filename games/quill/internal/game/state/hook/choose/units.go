package choose

import (
	"context"
	"slices"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const UnitsChoice = "Units"

type UnitsArgs struct {
	Types []string
}

func RetrieveUnits(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*UnitsArgs)
	units := make([]uuid.UUID, 0)
	for _, tile := range state.Board.UUIDs {
		if tile.Unit != nil {
			unit := tile.Unit.(*cd.UnitCard)
			if len(r.Types) == 0 || slices.Contains(r.Types, unit.Type) {
				units = append(units, unit.UUID)
			}
		}
	}
	return units, nil
}
