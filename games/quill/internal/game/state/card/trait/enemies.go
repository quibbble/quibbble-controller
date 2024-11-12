package trait

import (
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const (
	EnemiesTrait = "Enemies"
)

type EnemiesArgs struct {
	ChooseUnits parse.Choose
	Trait       parse.Trait

	// DO NOT SET IN YAML - SET BY ENGINE
	// current units that have the trait applied
	Current []uuid.UUID
}
