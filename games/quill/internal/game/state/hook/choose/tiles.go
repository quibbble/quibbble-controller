package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const TilesChoice = "Tiles"

type TilesArgs struct {
	Empty bool
}

func RetrieveTiles(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*TilesArgs)
	tiles := make([]uuid.UUID, 0)
	for _, tile := range state.Board.UUIDs {
		if (tile.Unit == nil) == r.Empty {
			tiles = append(tiles, tile.UUID)
		}
	}
	return tiles, nil
}
