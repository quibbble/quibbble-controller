package damage

import (
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/maths"
)

func Battle(state *st.State, attacker, defender *cd.UnitCard) (int, int, error) {
	attackerDamage, err := Damage(defender, maths.MaxInt(attacker.Attack, 0), attacker.DamageType)
	if err != nil {
		return 0, 0, errors.Wrap(err)
	}

	// assassin trait check
	assassins := attacker.GetTraits(tr.AssassinTrait)
	if len(assassins) > 0 {
		_, aY, err := state.Board.GetUnitXY(attacker.UUID)
		if err != nil {
			return 0, 0, errors.Wrap(err)
		}
		_, dY, err := state.Board.GetUnitXY(defender.UUID)
		if err != nil {
			return 0, 0, errors.Wrap(err)
		}
		defenderSide := state.Board.Sides[defender.Player]
		if maths.AbsInt(defenderSide-aY) < maths.AbsInt(defenderSide-dY) {
			for _, trait := range assassins {
				args := trait.GetArgs().(*tr.AssassinArgs)
				attackerDamage += args.Amount
			}
		}
	}

	defenderDamage, err := Damage(attacker, maths.MaxInt(defender.Attack, 0), defender.DamageType)
	if err != nil {
		return 0, 0, errors.Wrap(err)
	}

	// dodge trait check
	if len(defender.GetTraits(tr.DodgeTrait)) > 0 && state.Rand.Intn(3) == 0 {
		return 0, 0, nil
	}

	// spiky trait check
	for _, trait := range defender.GetTraits(tr.SpikyTrait) {
		args := trait.GetArgs().(*tr.SpikyArgs)
		defenderDamage += args.Amount
	}

	return attackerDamage, defenderDamage, nil
}
