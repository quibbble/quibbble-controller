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
	RecallUnitEvent = "RecallUnit"
)

type RecallUnitArgs struct {
	ChooseUnit parse.Choose

	// DO NOT SET IN YAML - SET BY ENGINE
	// tile location unit before recall
	ChooseTile parse.Choose
}

func RecallUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*RecallUnitArgs)
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)

	if unit.GetID() == "U0001" {
		return errors.Errorf("cannot rescind U0001")
	}

	state.Board.XYs[x][y].Unit = nil
	a.ChooseTile = parse.Choose{
		Type: ch.UUIDChoice,
		Args: ch.UUIDArgs{
			UUID: state.Board.XYs[x][y].UUID,
		},
	}

	// friends/enemies trait check
	for _, trait := range unit.GetTraits(tr.FriendsTrait) {
		args := trait.GetArgs().(*tr.FriendsArgs)
		if err := updateUnits(e, engine, state, unit.GetUUID(), args.Current, make([]uuid.UUID, 0), args.Trait); err != nil {
			return errors.Wrap(err)
		}
	}
	for _, trait := range unit.GetTraits(tr.EnemiesTrait) {
		args := trait.GetArgs().(*tr.EnemiesArgs)
		if err := updateUnits(e, engine, state, unit.GetUUID(), args.Current, make([]uuid.UUID, 0), args.Trait); err != nil {
			return errors.Wrap(err)
		}
	}
	if err := friendsTraitCheck(e, engine, state); err != nil {
		return errors.Wrap(err)
	}
	if err := enemiesTraitCheck(e, engine, state); err != nil {
		return errors.Wrap(err)
	}

	// reset and move items and unit back to hand
	for _, item := range unit.Items {
		item.Reset(state.BuildCard)
		state.Discard[item.Player].Add(item)
	}
	unit.Reset(state.BuildCard)

	if state.Hand[unit.Player].GetSize() > st.MaxHandSize {
		// burn the card if hand size to large
		state.Deck[unit.Player].Add(unit)
		event, err := NewEvent(state.Gen.New(en.EventUUID), BurnCardEvent, BurnCardArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: unit.Player,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(ctx, event, state); err != nil {
			return errors.Wrap(err)
		}
	} else {
		state.Hand[unit.Player].Add(unit)
	}
	return nil
}
