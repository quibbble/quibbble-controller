package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	dg "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/damage"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const (
	DamageUnitEvent = "DamageUnit"
)

type DamageUnitArgs struct {
	DamageType string
	Amount     int
	ChooseUnit parse.Choose

	// INTERNAL USE ONLY - do not redo damage calculations if damage is from a battle
	fromBattle bool
}

func DamageUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*DamageUnitArgs)
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)
	damage := a.Amount
	if !a.fromBattle {
		damage, err = dg.Damage(unit, a.Amount, a.DamageType)
		if err != nil {
			return errors.Wrap(err)
		}
	}
	unit.Health -= damage

	if unit.Health <= 0 {

		// eternal trait check
		for _, item := range unit.Items {
			eternals := item.GetTraits(tr.EternalTrait)
			if len(eternals) > 1 {
				return errors.Errorf("'%s' may only have one eternal trait", item.GetUUID())
			} else if len(eternals) == 1 {
				args := eternals[0].GetArgs().(*tr.EternalArgs)
				var conditions en.Conditions
				for _, c := range args.Conditions {
					condition, err := state.BuildCondition(state.Gen.New(en.ConditionUUID), c.Type, c.Not, c.Args)
					if err != nil {
						return errors.Wrap(err)
					}
					conditions = append(conditions, condition)
				}
				ctx := context.WithValue(context.WithValue(context.Background(), en.TraitEventCtx, e), en.TraitCardCtx, item.GetUUID())
				pass, err := conditions.Pass(ctx, engine, state)
				if err != nil {
					return errors.Wrap(err)
				}
				if !pass {
					continue
				}
				choices, err := ch.GetChoices(ctx, args.ChooseUnit, engine, state)
				if err != nil {
					return errors.Wrap(err)
				}
				if len(choices) == 0 {
					continue
				}
				if len(choices) != 1 {
					return errors.Errorf("eternal trait must retrieve one choice")
				}
				x, y, err := state.Board.GetUnitXY(choices[0])
				if err != nil {
					return errors.Wrap(err)
				}
				if err := unit.RemoveItem(item.GetUUID()); err != nil {
					return errors.Wrap(err)
				}
				nextUnit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)
				state.Hand[nextUnit.Player].Add(item)
				event, err := NewEvent(state.Gen.New(en.EventUUID), AddItemToUnitEvent, AddItemToUnitArgs{
					ChoosePlayer: parse.Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: nextUnit.Player,
						},
					},
					ChooseItem: parse.Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: item.GetUUID(),
						},
					},
					ChooseUnit: parse.Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: nextUnit.GetUUID(),
						},
					},
				})
				if err != nil {
					return errors.Wrap(err)
				}
				if err := engine.Do(ctx, event, state); err != nil {
					return errors.Wrap(err)
				}
			}
		}

		// kill unit if health <= 0
		event, err := NewEvent(state.Gen.New(en.EventUUID), KillUnitEvent, KillUnitArgs{
			ChooseUnit: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: unit.UUID,
				},
			},
			ChooseTile: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: state.Board.XYs[x][y].UUID,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	} else {
		// enrage trait check
		for _, trait := range unit.GetTraits(tr.EnrageTrait) {
			args := trait.GetArgs().(*tr.EnrageArgs)
			for _, h := range args.Hooks {
				hook, err := state.NewHook(state.Gen, unit.GetUUID(), h)
				if err != nil {
					return errors.Wrap(err)
				}
				engine.Register(hook)
			}
			for _, e := range args.Events {
				event, err := NewEvent(state.Gen.New(en.EventUUID), e.Type, e.Args)
				if err != nil {
					return errors.Wrap(err)
				}
				ctx := context.WithValue(context.WithValue(context.Background(), en.TraitCardCtx, unit.GetUUID()), en.TraitEventCtx, e)
				if err := engine.Do(ctx, event, state); err != nil {
					return errors.Wrap(err)
				}
			}
		}
	}

	return nil
}
