package stratego

import (
	"math"
	"slices"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type AI struct{}

var points = map[string]float64{
	flag:       100,
	bomb:       3,
	spy:        1,
	scout:      2,
	miner:      3,
	sergeant:   4,
	lieutenant: 5,
	captain:    6,
	major:      7,
	colonel:    8,
	general:    9,
	marshal:    10,
}

func (ai AI) Score(game qg.Game, team string) float64 {
	g := game.(*Stratego)
	if slices.Contains(g.winners, team) {
		return math.Inf(-1)
	}
	if len(g.winners) > 0 {
		return math.Inf(1)
	}
	score := 0.0

	for _, row := range g.board.board {
		for _, unit := range row {
			if unit.Team != nil && *unit.Team == team {
				score -= points[unit.Type]
			}
			if unit.Team != nil && *unit.Team != team {
				score += points[unit.Type]
			}
		}
	}
	return score
}
