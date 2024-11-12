package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const RemoveTraitFromCard = "RemoveTraitFromCard"

type RemoveTraitFromCardArgs struct {
	ChooseTrait parse.Choose
	ChooseCard  parse.Choose
}

func RemoveTraitFromCardAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*RemoveTraitFromCardArgs)
	traitChoice, err := ch.GetChoice(ctx, a.ChooseTrait, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	cardChoice, err := ch.GetChoice(ctx, a.ChooseCard, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card := state.GetCard(cardChoice)
	if card == nil {
		return st.ErrNotFound(cardChoice)
	}
	if err := card.RemoveTrait(traitChoice); err != nil {
		return errors.Wrap(err)
	}

	// check if killed
	if _, _, err := state.Board.GetUnitXY(card.GetUUID()); err == nil && card.(*cd.UnitCard).Health <= 0 {
		event, err := NewEvent(state.Gen.New(en.EventUUID), KillUnitEvent, KillUnitArgs{
			ChooseUnit: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: card.GetUUID(),
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}

	// friends/enemies trait check
	if err := friendsTraitCheck(e, engine, state); err != nil {
		return errors.Wrap(err)
	}
	if err := enemiesTraitCheck(e, engine, state); err != nil {
		return errors.Wrap(err)
	}

	return nil
}
