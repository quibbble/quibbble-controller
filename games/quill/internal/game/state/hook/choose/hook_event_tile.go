package choose

import (
	"context"

	"github.com/go-viper/mapstructure/v2"
	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const HookEventTileChoice = "HookEventTile"

type HookEventTileArgs struct{}

func RetrieveHookTileUnit(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {

	event := ctx.Value(en.HookEventCtx).(en.IEvent)

	var a struct {
		ChooseTile parse.Choose
	}
	if err := mapstructure.Decode(event.GetArgs(), &a); err != nil {
		return nil, errors.ErrInterfaceConversion
	}

	tileChoice, err := GetUnitChoice(ctx, a.ChooseTile, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return []uuid.UUID{tileChoice}, nil
}
