package damage

import (
	"slices"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const (
	PureDamage     = "Pure"
	PhysicalDamage = "Physical"
	RangedDamage   = "Ranged" // SubType of PhysicalDamage
	MagicDamage    = "Magic"
	PoisonDamage   = "Poison" // SubType of MagicDamage
)

func Damage(unit *cd.UnitCard, amount int, typ string) (int, error) {
	if !slices.Contains([]string{PureDamage, PhysicalDamage, MagicDamage, RangedDamage, PoisonDamage}, typ) {
		return 0, errors.Errorf("'%s' is not a valid damage type", typ)
	}
	reduction := 0
	for _, trait := range unit.Traits {
		if slices.Contains([]string{PhysicalDamage, RangedDamage}, typ) && trait.GetType() == tr.ShieldTrait {
			args := trait.GetArgs().(*tr.ShieldArgs)
			reduction += args.Amount
		} else if slices.Contains([]string{MagicDamage, PoisonDamage}, typ) && trait.GetType() == tr.WardTrait {
			args := trait.GetArgs().(*tr.WardArgs)
			reduction += args.Amount
		}
	}
	damage := amount - reduction
	if damage < 0 {
		damage = 0
	}
	return damage, nil
}
