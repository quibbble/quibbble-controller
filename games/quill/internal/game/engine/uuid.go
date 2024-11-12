package engine

import (
	"strings"

	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const (
	TileUUID      = 'T'
	UnitUUID      = 'U'
	SpellUUID     = 'S'
	ItemUUID      = 'I'
	PlayerUUID    = 'P'
	EngineUUID    = 'E'
	EventUUID     = 'V'
	TraitUUID     = 'R'
	HookUUID      = 'H'
	ChooseUUID    = 'O'
	ConditionUUID = 'C'
	SackUUID      = 'K'
)

var (
	ErrInvalidUUIDType = func(uuid uuid.UUID, expectedTypes ...rune) error {
		str := make([]string, 0)
		for _, typ := range expectedTypes {
			str = append(str, string(typ))
		}
		return errors.Errorf("'%s' is not of type '%s'", string(uuid), strings.Join(str, ", "))
	}
)
