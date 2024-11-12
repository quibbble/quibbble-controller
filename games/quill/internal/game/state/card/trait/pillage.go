package trait

import "github.com/quibbble/quibbble-controller/games/quill/parse"

const (
	PillageTrait = "Pillage"
)

type PillageArgs struct {
	Description string
	Hooks       []parse.Hook
	Events      []parse.Event
}
