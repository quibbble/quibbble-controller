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
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const (
	AddTraitToCard = "AddTraitToCard"
)

type AddTraitToCardArgs struct {
	Trait      parse.Trait
	ChooseCard parse.Choose

	// NOT SET IN YAML - SET BY ENGINE
	// which item/spell/unit created the trait
	createdBy *uuid.UUID
}

func AddTraitToCardAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*AddTraitToCardArgs)
	trait, err := tr.NewTrait(state.Gen.New(en.ChooseUUID), a.createdBy, a.Trait.Type, a.Trait.Args)
	if err != nil {
		return errors.Wrap(err)
	}

	choice, err := ch.GetChoice(ctx, a.ChooseCard, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card := state.GetCard(choice)
	if card == nil {
		return st.ErrNotFound(choice)
	}
	if err := card.AddTrait(trait); err != nil {
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
