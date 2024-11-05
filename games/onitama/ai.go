package onitama

import (
	"math"
	"slices"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type AI struct{}

func (ai AI) Score(game qg.Game, team string) float64 {
	g := game.(*Onitama)
	if slices.Contains(g.winners, team) {
		return math.Inf(1)
	}
	if len(g.winners) > 0 {
		return math.Inf(-1)
	}
	myPawns := 0
	enemyPawns := 0
	for _, row := range g.board {
		for _, pawn := range row {
			if pawn != nil && pawn.Team == team {
				myPawns++
			} else if pawn != nil && pawn.Team != team {
				enemyPawns++
			}
		}
	}
	return float64(myPawns) + (5.0 - float64(enemyPawns))
}
