package game

import (
	"fmt"
	"math"
	"strconv"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

// Actions common to all games.
const (
	UndoAction  = "undo"
	ResetAction = "reset"
	AIAction    = "ai"
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

// AI plays the turn for the current player.
func AI(b GameBuilder, ai GameAI, g Game, depth int) error {
	snapshot, err := g.GetSnapshotJSON()
	if err != nil {
		return err
	}
	if len(snapshot.Winners) > 0 {
		return fmt.Errorf("game already over")
	}
	_, action, err := alphabeta(b, ai, g, depth, math.Inf(-1), math.Inf(1), snapshot.Turn)
	if err != nil {
		return err
	}
	return g.Do(action)
}
