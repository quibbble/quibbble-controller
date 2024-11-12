package trait

import "github.com/quibbble/quibbble-controller/games/quill/parse"

const (
	EnrageTrait = "Enrage"
)

type EnrageArgs struct {
	Description string
	Hooks       []parse.Hook
	Events      []parse.Event
}
