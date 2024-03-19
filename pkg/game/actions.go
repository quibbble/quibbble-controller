package game

import (
	"strconv"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

// Actions common to all games.
const (
	UndoAction  = "undo"
	ResetAction = "reset"
)

// Undo returns a new game with the last action  undone.
func Undo(b GameBuilder, g Game) (Game, error) {
	snapshot, err := g.GetSnapshotQGN()
	if err != nil {
		return nil, err
	}
	if len(snapshot.Actions) > 0 {
		snapshot.Actions = snapshot.Actions[:len(snapshot.Actions)-1]
	}
	return b.Create(snapshot)
}

// Reset resets the game to base state. If seed
// is passed and the game utilizes randomness then
// the new game will use this new seed.
func Reset(b GameBuilder, g Game, seed ...int) (Game, error) {
	snapshot, err := g.GetSnapshotQGN()
	if err != nil {
		return nil, err
	}
	if len(seed) > 0 {
		snapshot.Tags[qgn.SeedTag] = strconv.Itoa(seed[0])
	}
	snapshot.Actions = make([]qgn.Action, 0)
	return b.Create(snapshot)
}
