package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const (
	DrainManaEvent = "DrainMana"
)

type DrainManaArgs struct {
	ChoosePlayer parse.Choose
	Amount       int
}

func DrainManaAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*DrainManaArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	state.Mana[playerChoice].Amount -= a.Amount
	if state.Mana[playerChoice].Amount < 0 {
		state.Mana[playerChoice].Amount = 0
	}

	// surge trait check
	for _, col := range state.Board.XYs {
		for _, tile := range col {
			if tile.Unit != nil && tile.Unit.GetPlayer() == playerChoice {
				for range tile.Unit.GetTraits(tr.SurgeTrait) {
					tile.Unit.(*cd.UnitCard).Attack -= a.Amount
				}
			}
		}
	}
	return nil
}
