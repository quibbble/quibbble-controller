package trait

import (
	"github.com/quibbble/quibbble-controller/games/quill/parse"
)

const (
	BattleCryTrait = "BattleCry"
)

type BattleCryArgs struct {
	Description string
	Hooks       []parse.Hook
	Events      []parse.Event
}
