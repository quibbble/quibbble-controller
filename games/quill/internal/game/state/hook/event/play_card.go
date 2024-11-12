package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/maths"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const (
	PlayCardEvent = "PlayCard"
)

type PlayCardArgs struct {
	ChoosePlayer parse.Choose
	ChooseCard   parse.Choose
}

func PlayCardAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*PlayCardArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	cardChoice, err := ch.GetChoice(ctx, a.ChooseCard, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card, err := state.Hand[playerChoice].GetCard(cardChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	playable, err := card.Playable(engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	next, err := card.NextTargets(ctx, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	if (len(card.GetTargets()) != len(ctx.Value(en.TargetsCtx).([]uuid.UUID))) || !playable || len(next) != 0 {
		return errors.Errorf("card cannot be played")
	}

	// drain mana equal to card cost
	if state.Mana[playerChoice].Amount < card.GetCost() {
		return errors.Errorf("player '%s' does not have enough mana to play '%s'", playerChoice, card.GetUUID())
	}

	targets := ctx.Value(en.TargetsCtx).([]uuid.UUID)

	// purity trait check
	if card.GetUUID().Type() == en.SpellUUID {
		for _, target := range targets {
			if x, y, err := state.Board.GetUnitXY(target); err == nil {
				unit := state.Board.XYs[x][y].Unit
				if len(unit.GetTraits(tr.PurityTrait)) > 0 {
					return errors.Errorf("'%s' cannot target '%s' due to purity trait", card.GetUUID(), unit.GetUUID())
				}
			}

		}
	}

	event, err := NewEvent(state.Gen.New(en.EventUUID), DrainManaEvent, DrainManaArgs{
		ChoosePlayer: parse.Choose{
			Type: ch.CurrentPlayerChoice,
			Args: ch.CurrentPlayerArgs{},
		},
		Amount: maths.MaxInt(card.GetCost(), 0),
	})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := engine.Do(context.Background(), event, state); err != nil {
		return errors.Wrap(err)
	}

	// add hooks
	for _, hook := range card.GetHooks() {
		engine.Register(hook)
	}

	// create event based on card type
	switch card.GetUUID().Type() {
	case en.ItemUUID:
		if len(targets) <= 0 {
			return errors.ErrIndexOutOfBounds
		}
		event, err = NewEvent(state.Gen.New(en.EventUUID), AddItemToUnitEvent, AddItemToUnitArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: playerChoice,
				},
			},
			ChooseItem: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: card.GetUUID(),
				},
			},
			ChooseUnit: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: targets[0],
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
	case en.SpellUUID:
		event, err = NewEvent(state.Gen.New(en.EventUUID), DiscardCardEvent, DiscardCardArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.CurrentPlayerChoice,
				Args: ch.CurrentPlayerArgs{},
			},
			ChooseCard: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: card.GetUUID(),
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
	case en.UnitUUID:
		if len(targets) <= 0 {
			return errors.ErrIndexOutOfBounds
		}
		event, err = NewEvent(state.Gen.New(en.EventUUID), PlaceUnitEvent, PlaceUnitArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.CurrentPlayerChoice,
				Args: ch.CurrentPlayerArgs{},
			},
			ChooseUnit: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: card.GetUUID(),
				},
			},
			ChooseTile: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: targets[0],
				},
			},
			InPlayRange: true,
		})
		if err != nil {
			return errors.Wrap(err)
		}
	default:
		return errors.ErrMissingMapKey
	}

	// apply card event and then any additional events attached to the card
	events := append([]en.IEvent{event}, card.GetEvents()...)
	for _, event := range events {
		if err := engine.Do(ctx, event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
