package trait

import "github.com/quibbble/quibbble-controller/games/quill/parse"

const (
	EternalTrait = "Eternal"
)

type EternalArgs struct {
	Conditions []parse.Condition
	ChooseUnit parse.Choose
}
