package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const UUIDChoice = "UUID"

type UUIDArgs struct {
	UUID uuid.UUID
}

func RetrieveUUID(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*UUIDArgs)
	return []uuid.UUID{r.UUID}, nil
}
