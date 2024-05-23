package indigo

import (
	"math"
	"slices"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type AI struct{}

func (ai AI) Score(game qg.Game, team string) float64 {
	g := game.(*Indigo)
	if slices.Contains(g.winners, team) {
		return math.Inf(1)
	}
	if len(g.winners) > 0 {
		return math.Inf(-1)
	}

	max := 0
	for t, p := range g.points {
		if t != team {
			if p > max {
				max = p
			}
		}
	}
	return float64(max - g.points[team])
}
