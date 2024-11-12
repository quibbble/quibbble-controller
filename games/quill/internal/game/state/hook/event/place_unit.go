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
)

const (
	PlaceUnitEvent = "PlaceUnit"
)

type PlaceUnitArgs struct {
	ChoosePlayer parse.Choose
	ChooseUnit   parse.Choose
	ChooseTile   parse.Choose
	InPlayRange  bool
}

func PlaceUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*PlaceUnitArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	tileChoice, err := ch.GetTileChoice(ctx, a.ChooseTile, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card, err := state.Hand[playerChoice].GetCard(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unit := card.(*cd.UnitCard)

	tX, tY, err := state.Board.GetTileXY(tileChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	if state.Board.XYs[tX][tY].Unit != nil {
		return errors.Errorf("unit '%s' cannot be placed on a full tile", unit.UUID)
	}
	if a.InPlayRange {
		min, max := state.Board.GetPlayableRowRange(playerChoice)
		if tY < min || tY > max {
			return errors.Errorf("unit '%s' must be placed within rows %d to %d", unit.UUID, min, max)
		}
	}
	if err := state.Hand[playerChoice].RemoveCard(unitChoice); err != nil {
		return errors.Wrap(err)
	}

	// haste trait check
	if len(unit.GetTraits(tr.HasteTrait)) > 0 {
		unit.Cooldown = 0
	}

	// surge trait check
	for range unit.GetTraits(tr.SurgeTrait) {
		unit.Attack += state.Mana[playerChoice].Amount
	}

	state.Board.XYs[tX][tY].Unit = unit

	// battle cry trait check
	for _, trait := range unit.GetTraits(tr.BattleCryTrait) {
		args := trait.GetArgs().(*tr.BattleCryArgs)
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

	// friends/enemies trait check
	if err := friendsTraitCheck(e, engine, state); err != nil {
		return errors.Wrap(err)
	}
	if err := enemiesTraitCheck(e, engine, state); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
