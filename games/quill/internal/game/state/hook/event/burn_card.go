package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const (
	BurnCardEvent = "BurnCard"
)

type BurnCardArgs struct {
	ChoosePlayer parse.Choose
}

func BurnCardAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*BurnCardArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	card, err := state.Deck[playerChoice].Draw()
	if err != nil {
		return errors.Wrap(err)
	}
	state.Trash[playerChoice].Add(*card)
	return nil
}
