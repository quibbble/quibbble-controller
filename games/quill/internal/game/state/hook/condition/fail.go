package condition

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
)

const FailCondition = "Fail"

type FailArgs struct{}

func PassFail(c *Condition, ctx context.Context, engine *en.Engine, state *st.State) (bool, error) {
	return false, nil
}
