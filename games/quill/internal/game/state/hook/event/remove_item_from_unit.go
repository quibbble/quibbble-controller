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
	RemoveItemFromUnitEvent = "RemoveItemFromUnit"
)

type RemoveItemFromUnitArgs struct {
	ChooseItem parse.Choose
	ChooseUnit parse.Choose
}

func RemoveItemFromUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*RemoveItemFromUnitArgs)
	itemChoice, err := ch.GetItemChoice(ctx, a.ChooseItem, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	card := state.Board.XYs[x][y].Unit.(*cd.UnitCard)
	item, err := card.GetAndRemoveItem(itemChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, trait := range card.Traits {
		createdBy := trait.GetCreatedBy()
		if createdBy != nil && *createdBy == item.GetUUID() {
			event, err := NewEvent(state.Gen.New(en.EventUUID), RemoveTraitFromCard, RemoveTraitFromCardArgs{
				ChooseTrait: parse.Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: trait.GetUUID(),
					},
				},
				ChooseCard: parse.Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: unitChoice,
					},
				},
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
	}
	return nil
}
