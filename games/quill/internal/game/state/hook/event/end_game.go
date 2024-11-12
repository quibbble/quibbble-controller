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
	EndGameEvent = "EndGame"
)

type EndGameArgs struct {
	ChooseWinner parse.Choose
}

func EndGameAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*EndGameArgs)
	player, err := ch.GetPlayerChoice(ctx, a.ChooseWinner, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	state.Winner = &player
	return nil
}
