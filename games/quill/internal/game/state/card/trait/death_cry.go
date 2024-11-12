package trait

import "github.com/quibbble/quibbble-controller/games/quill/parse"

const (
	DeathCryTrait = "DeathCry"
)

type DeathCryArgs struct {
	Description string
	Hooks       []parse.Hook
	Events      []parse.Event
}
