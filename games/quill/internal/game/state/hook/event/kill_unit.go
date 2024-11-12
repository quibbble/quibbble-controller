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
	KillUnitEvent = "KillUnit"
)

type KillUnitArgs struct {
	ChooseUnit parse.Choose

	// DO NOT SET IN YAML - SET BY ENGINE
	// tile location of unit before death
	ChooseTile parse.Choose
}

func KillUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*KillUnitArgs)
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)

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

	if unit.GetID() != "U0001" {
		// reset and move items and unit to discard
		for _, item := range unit.Items {
			if !item.Token {
				item.Reset(state.BuildCard)
				state.Discard[item.Player].Add(item)
			}
		}
		if !unit.Token {
			unit.Reset(state.BuildCard)
			state.Discard[unit.Player].Add(unit)
		}
	} else {
		// check if the game is over
		choose1, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), ch.UnitsChoice, ch.UnitsArgs{
			Types: []string{cd.BaseUnit},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		choose2, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), ch.OwnedUnitsChoice, ch.OwnedUnitsArgs{
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
		choices := ch.NewChooseChain(ch.SetIntersect, choose1, choose2)
		bases, err := choices.Retrieve(context.Background(), engine, state)
		if err != nil {
			return errors.Wrap(err)
		}
		if len(bases) <= st.Cols/2 {
			event, err := NewEvent(state.Gen.New(en.EventUUID), EndGameEvent, EndGameArgs{
				ChooseWinner: parse.Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: state.GetOpponent(unit.Player),
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
	}

	// death cry trait check
	for _, trait := range unit.GetTraits(tr.DeathCryTrait) {
		args := trait.GetArgs().(*tr.DeathCryArgs)
		for _, h := range args.Hooks {
			hook, err := state.NewHook(state.Gen, unit.GetUUID(), h)
			if err != nil {
				return errors.Wrap(err)
			}
			engine.Register(hook)
		}
		for _, ev := range args.Events {
			event, err := NewEvent(state.Gen.New(en.EventUUID), ev.Type, ev.Args)
			if err != nil {
				return errors.Wrap(err)
			}
			ctx := context.WithValue(context.WithValue(context.Background(), en.TraitCardCtx, unit.GetUUID()), en.TraitEventCtx, e)
			if err := engine.Do(ctx, event, state); err != nil {
				return errors.Wrap(err)
			}
		}
	}
	return nil
}
