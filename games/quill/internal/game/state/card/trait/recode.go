package trait

import (
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"

	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const RecodeTrait = "Recode"

const (
	Union     = "Union"
	Intersect = "Intersect"
	Replace   = "Replace"
)

type RecodeArgs struct {
	Code string

	SetFunction string // Union, Intersect, Replace

	// INTERNAL USE ONLY
	oldCode string
}

func AddRecode(t *Trait, card st.ICard) error {
	a := t.GetArgs().(*RecodeArgs)
	unit, ok := card.(*cd.UnitCard)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	switch a.SetFunction {
	case Union:
		code := ""
		for i, c := range a.Code {
			if unit.Codex[i] == '1' {
				code += "1"
			} else {
				code += string(c)
			}
		}
		unit.Codex = code
	case Intersect:
		code := ""
		for i, c := range a.Code {
			if unit.Codex[i] == '1' && c == '1' {
				code += "1"
			} else {
				code += "0"
			}
		}
		unit.Codex = code
	case Replace:
		unit.Codex = a.Code
	default:
		return errors.Errorf("invalid set function %s", a.SetFunction)
	}
	return nil
}

func RemoveRecode(t *Trait, card st.ICard) error {
	a := t.GetArgs().(*RecodeArgs)
	unit, ok := card.(*cd.UnitCard)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	unit.Codex = a.oldCode
	return nil
}
