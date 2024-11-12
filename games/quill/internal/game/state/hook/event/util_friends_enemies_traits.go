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

func friendsTraitCheck(e *Event, engine *en.Engine, state *st.State) error {
	// checks if any units need to add/remove their trait through friend trait
	for _, col := range state.Board.XYs {
		for _, tile := range col {
			if tile.Unit != nil {
				for _, trait := range tile.Unit.GetTraits(tr.FriendsTrait) {
					args := trait.GetArgs().(*tr.FriendsArgs)
					before := args.Current
					if before == nil {
						before = make([]uuid.UUID, 0)
					}
					choose1, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), args.ChooseUnits.Type, args.ChooseUnits.Args)
					if err != nil {
						return errors.Wrap(err)
					}
					choose2, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), ch.OwnedUnitsChoice, &ch.OwnedUnitsArgs{
						ChoosePlayer: parse.Choose{
							Type: ch.UUIDChoice,
							Args: ch.UUIDArgs{
								UUID: tile.Unit.GetPlayer(),
							},
						},
					})
					if err != nil {
						return errors.Wrap(err)
					}
					after, err := ch.NewChooseChain(ch.SetIntersect, choose1, choose2).Retrieve(context.WithValue(context.Background(), en.CardCtx, tile.Unit.GetUUID()), engine, state)
					if err != nil {
						return errors.Wrap(err)
					}
					args.Current = after
					if err := updateUnits(e, engine, state, tile.Unit.GetUUID(), before, after, args.Trait); err != nil {
						return errors.Wrap(err)
					}
				}
			}
		}
	}
	return nil
}

func enemiesTraitCheck(e *Event, engine *en.Engine, state *st.State) error {
	// checks if any units need to add/remove their trait through friend trait
	for _, col := range state.Board.XYs {
		for _, tile := range col {
			if tile.Unit != nil {
				for _, trait := range tile.Unit.GetTraits(tr.EnemiesTrait) {
					args := trait.GetArgs().(*tr.EnemiesArgs)
					before := args.Current
					if before == nil {
						before = make([]uuid.UUID, 0)
					}
					choose1, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), args.ChooseUnits.Type, args.ChooseUnits.Args)
					if err != nil {
						return errors.Wrap(err)
					}
					choose2, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), ch.OwnedUnitsChoice, &ch.OwnedUnitsArgs{
						ChoosePlayer: parse.Choose{
							Type: ch.UUIDChoice,
							Args: ch.UUIDArgs{
								UUID: state.GetOpponent(tile.Unit.GetPlayer()),
							},
						},
					})
					if err != nil {
						return errors.Wrap(err)
					}
					after, err := ch.NewChooseChain(ch.SetIntersect, choose1, choose2).Retrieve(context.WithValue(context.Background(), en.CardCtx, tile.Unit.GetUUID()), engine, state)
					if err != nil {
						return errors.Wrap(err)
					}
					args.Current = after
					if err := updateUnits(e, engine, state, tile.Unit.GetUUID(), before, after, args.Trait); err != nil {
						return errors.Wrap(err)
					}
				}
			}
		}
	}
	return nil
}

func updateUnits(e *Event, engine *en.Engine, state *st.State, createdBy uuid.UUID, before, after []uuid.UUID, trait parse.Trait) error {
	remove := uuid.Diff(before, after)
	for _, u := range remove {
		x, y, err := state.Board.GetUnitXY(u)
		if err != nil {
			return errors.Wrap(err)
		}
		found := false
		unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)
		for _, t := range unit.GetTraits(trait.Type) {
			if t.GetCreatedBy() != nil && *t.GetCreatedBy() == createdBy {
				// if reflect.DeepEqual(t.GetArgs(), trait.GetArgs()) {
				event, err := NewEvent(state.Gen.New(en.EventUUID), RemoveTraitFromCard, &RemoveTraitFromCardArgs{
					ChooseTrait: parse.Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: t.GetUUID(),
						},
					},
					ChooseCard: parse.Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: u,
						},
					},
				})
				if err != nil {
					return errors.Wrap(err)
				}
				ctx := context.WithValue(context.WithValue(context.Background(), en.TraitCardCtx, unit.GetUUID()), en.TraitEventCtx, e)
				if err := engine.Do(ctx, event, state); err != nil {
					return errors.Wrap(err)
				}
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf("failed to find '%s' trait for '%s'", trait.Type, u)
		}
	}
	add := uuid.Diff(after, before)
	for _, u := range add {
		event, err := NewEvent(state.Gen.New(en.EventUUID), AddTraitToCard, &AddTraitToCardArgs{
			Trait: parse.Trait{
				Type: trait.Type,
				Args: trait.Args,
			},
			ChooseCard: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: u,
				},
			},
			createdBy: &createdBy,
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
