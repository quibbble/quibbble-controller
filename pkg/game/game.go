package game

import qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"

type Game interface {
	// Do performs an action on the game
	Do(action *Action) error

	// GetSnapshotQGN returns the current game state in Quibbble Game Notation
	GetSnapshotQGN() (*qgn.Snapshot, error)

	// GetSnapshotJSON returns the current game state from 'team' view
	// Entering nothing returns a complete snapshot with no data hidden i.e. all hands, resources, etc.
	// Entering more than one team should error.
	GetSnapshotJSON(team ...string) (*Snapshot, error)
}
