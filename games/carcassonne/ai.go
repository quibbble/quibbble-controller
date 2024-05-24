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

	return float64(g.scores[team])
}
