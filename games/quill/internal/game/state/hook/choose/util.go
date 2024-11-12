package choose

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

func GetPlayerChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	return getTypeChoice(ctx, raw, en.PlayerUUID, engine, state)
}

func GetUnitChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	return getTypeChoice(ctx, raw, en.UnitUUID, engine, state)
}

func GetItemChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	return getTypeChoice(ctx, raw, en.ItemUUID, engine, state)
}

func GetSpellChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	return getTypeChoice(ctx, raw, en.SpellUUID, engine, state)
}

func GetTileChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	return getTypeChoice(ctx, raw, en.TileUUID, engine, state)
}

func GetChoices(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	choose, err := NewChoose(state.Gen.New(en.ChooseUUID), raw.Type, raw.Args)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return choices, nil
}

func GetChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	choices, err := GetChoices(ctx, raw, engine, state)
	if err != nil {
		return uuid.Nil, errors.Wrap(err)
	}
	if len(choices) != 1 {
		return uuid.Nil, errors.ErrInvalidSliceLength
	}
	return choices[0], nil
}

func getTypeChoice(ctx context.Context, raw parse.Choose, typ rune, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	choice, err := GetChoice(ctx, raw, engine, state)
	if err != nil {
		return uuid.Nil, errors.Wrap(err)
	}
	if choice.Type() != typ {
		return uuid.Nil, en.ErrInvalidUUIDType(choice, typ)
	}
	return choice, nil
}
