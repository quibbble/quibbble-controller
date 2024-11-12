package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const TargetChoice = "Target"

type TargetArgs struct {
	Index int
}

func RetrieveTarget(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*TargetArgs)
	targets := ctx.Value(en.TargetsCtx).([]uuid.UUID)
	if r.Index < 0 || r.Index >= len(targets) {
		return nil, errors.ErrIndexOutOfBounds
	}
	return []uuid.UUID{targets[r.Index]}, nil
}
