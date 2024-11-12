package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const DrawCardEvent = "DrawCard"

type DrawCardArgs struct {
	ChoosePlayer parse.Choose
}

func DrawCardAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*DrawCardArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	if state.Hand[playerChoice].GetSize() <= st.MaxHandSize {
		card, err := state.Deck[playerChoice].Draw()
		if err != nil {
			return errors.Wrap(err)
		}
		state.Hand[playerChoice].Add(*card)
	} else {
		event, err := NewEvent(state.Gen.New(en.EventUUID), BurnCardEvent, BurnCardArgs{
			ChoosePlayer: a.ChoosePlayer,
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
