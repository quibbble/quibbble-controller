package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const CardIDByCostChoice = "CardIDByCost"

type CardIDByCostArgs struct {
	Cost int
}

func RetrieveCardIDByCost(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	args := c.Args.(*CardIDByCostArgs)
	choices := make([]uuid.UUID, 0)
	for _, id := range st.CardIDByCostMap[args.Cost] {
		choices = append(choices, uuid.UUID(id))
	}
	return choices, nil
}
