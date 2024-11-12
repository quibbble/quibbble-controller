package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const CardIDByTypeChoice = "CardIDByType"

type CardIDByTypeArgs struct {
	CardTypes []string
}

func RetrieveCardIDByType(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	args := c.Args.(*CardIDByTypeArgs)
	choices := make([]uuid.UUID, 0)
	for _, typ := range args.CardTypes {
		for _, id := range st.CardIDByTypeMap[typ] {
			choices = append(choices, uuid.UUID(id))
		}
	}
	return choices, nil
}
