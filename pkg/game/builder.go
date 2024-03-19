package game

import qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"

type GameBuilder interface {
	// Create creates/recreates a game using the given Quibbble Game Notation
	Create(game *qgn.Snapshot) (Game, error)

	// GetInformation provides additional details about the game
	GetInformation() *Information
}
