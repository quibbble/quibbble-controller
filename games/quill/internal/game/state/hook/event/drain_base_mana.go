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
	DrainBaseManaEvent = "DrainBaseMana"
)

type DrainBaseManaArgs struct {
	ChoosePlayer parse.Choose
	Amount       int
}

func DrainBaseManaAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*DrainBaseManaArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	state.Mana[playerChoice].BaseAmount -= a.Amount
	if state.Mana[playerChoice].BaseAmount < 0 {
		state.Mana[playerChoice].BaseAmount = 0
	}
	return nil
}
