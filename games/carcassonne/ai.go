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

	tileScore, err := g.score(team)
	if err != nil {
		return 0
	}

	score := tileScore
	if g.tokens[team] >= 5 {
		score -= float64(g.tokens[team]) // prioritize using more tokens when you have a lot
	} else if g.tokens[team] <= 2 {
		score += float64(g.tokens[team] * 2) // prioritize saving tokens if you only have a few
	}

	return score
}
