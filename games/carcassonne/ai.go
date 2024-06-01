package carcassonne

import (
	"math"
	"slices"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type AI struct{}

func (ai AI) Score(game qg.Game, team string) float64 {
	g := game.(*Carcassonne)
	if slices.Contains(g.winners, team) {
		return math.Inf(1)
	}
	if len(g.winners) > 0 {
		return math.Inf(-1)
	}

	// skip scoring after place tile action
	if g.playTiles[team] == nil {
		return float64(0)
	}

	score, err := g.scoreLastTile(team)
	if err != nil {
		return 0
	}

	return float64(score + g.tokens[team])
}
