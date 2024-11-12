package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const SelfChoice = "Self"

type SelfArgs struct{}

func RetrieveSelf(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	cardCtx := ctx.Value(en.HookCardCtx)
	if cardCtx == nil {
		cardCtx = ctx.Value(en.TraitCardCtx)
		if cardCtx == nil {
			cardCtx = ctx.Value(en.CardCtx)
			if cardCtx == nil {
				return nil, errors.ErrMissingContext
			}
		}
	}
	cardUUID := cardCtx.(uuid.UUID)
	return []uuid.UUID{cardUUID}, nil
}
