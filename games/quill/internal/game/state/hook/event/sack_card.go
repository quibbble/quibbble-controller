package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const SackCardEvent = "SackCard"

const (
	ManaSackOption  = "Mana"
	CardsSackOption = "Cards"
)

type SackCardArgs struct {
	ChoosePlayer parse.Choose
	SackOption   string
	ChooseCard   parse.Choose
}

func SackCardAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*SackCardArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	event, err := NewEvent(state.Gen.New(en.EventUUID), DiscardCardEvent, DiscardCardArgs{
		ChoosePlayer: parse.Choose{
			Type: ch.CurrentPlayerChoice,
			Args: ch.CurrentPlayerArgs{},
		},
		ChooseCard: a.ChooseCard,
	})
	if err != nil {
		return errors.Wrap(err)
	}

	events := []en.IEvent{event}

	switch a.SackOption {
	case ManaSackOption:
		event1, err := NewEvent(state.Gen.New(en.EventUUID), GainBaseManaEvent, GainBaseManaArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.CurrentPlayerChoice,
				Args: ch.CurrentPlayerArgs{},
			},
			Amount: 1,
		})
		if err != nil {
			return errors.Wrap(err)
		}
		event2, err := NewEvent(state.Gen.New(en.EventUUID), GainManaEvent, GainManaArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.CurrentPlayerChoice,
				Args: ch.CurrentPlayerArgs{},
			},
			Amount: 1,
		})
		if err != nil {
			return errors.Wrap(err)
		}
		events = append(events, event1, event2)
	case CardsSackOption:
		event1, err := NewEvent(state.Gen.New(en.EventUUID), DrawCardEvent, DrawCardArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: playerChoice,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		event2, err := NewEvent(state.Gen.New(en.EventUUID), DrawCardEvent, DrawCardArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: playerChoice,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		events = append(events, event1, event2)
	default:
		return errors.Errorf("invalid sack option '%s'", a.SackOption)
	}

	for _, event := range events {
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}

	state.Sacked[playerChoice] = true

	return nil
}
