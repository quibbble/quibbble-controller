package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const CurrentPlayerChoice = "CurrentPlayer"

type CurrentPlayerArgs struct{}

func RetrieveCurrentPlayer(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	return []uuid.UUID{state.GetTurn()}, nil
}
