package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const OpposingPlayerChoice = "OpposingPlayer"

type OpposingPlayerArgs struct{}

func RetrieveOpposingPlayer(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	return []uuid.UUID{state.GetOpponent(state.GetTurn())}, nil
}
