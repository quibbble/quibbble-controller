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

const (
	AddItemToUnitEvent = "AddItemToUnit"
)

type AddItemToUnitArgs struct {
	ChoosePlayer parse.Choose
	ChooseItem   parse.Choose
	ChooseUnit   parse.Choose
}

func AddItemToUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*AddItemToUnitArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	itemChoice, err := ch.GetItemChoice(ctx, a.ChooseItem, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card, err := state.Hand[playerChoice].GetCard(itemChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	itemCard := card.(*cd.ItemCard)

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unitCard := state.Board.XYs[x][y].Unit.(*cd.UnitCard)

	if unitCard.Type == cd.StructureUnit {
		return errors.Errorf("cannot add an item to a structure unit")
	}
	if err := state.Hand[playerChoice].RemoveCard(itemChoice); err != nil {
		return errors.Wrap(err)
	}
	if err := unitCard.AddItem(itemCard); err != nil {
		return errors.Wrap(err)
	}

	for _, trait := range itemCard.HeldTraits {
		event, err := NewEvent(state.Gen.New(en.EventUUID), AddTraitToCard, AddTraitToCardArgs{
			Trait: parse.Trait{
				Type: trait.GetType(),
				Args: trait.GetArgs(),
			},
			ChooseCard: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: unitChoice,
				},
			},
			createdBy: &itemChoice,
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}

		// if unit died from adding trait then break
		_, _, err = state.Board.GetUnitXY(unitChoice)
		if err != nil {
			break
		}
	}
	return nil
}
