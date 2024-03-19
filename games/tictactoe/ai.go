package tictactoe

import (
	"math"
	"slices"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type AI struct{}

func (ai AI) Score(game qg.Game, team string) float64 {
	g := game.(*Tictactoe)
	if slices.Contains(g.winners, team) {
		return math.Inf(-1)
	}
	if len(g.winners) > 0 {
		return math.Inf(1)
	}
	score := 0.0
	for i := range BoardSize {
		if (g.board[i][0] == nil || *g.board[i][0] == team) &&
			(g.board[i][1] == nil || *g.board[i][1] == team) &&
			(g.board[i][2] == nil || *g.board[i][2] == team) {
			score -= 10
		}
		if (g.board[0][i] == nil || *g.board[0][i] == team) &&
			(g.board[1][i] == nil || *g.board[1][i] == team) &&
			(g.board[2][i] == nil || *g.board[2][i] == team) {
			score -= 10
		}
	}
	if (g.board[0][0] == nil || *g.board[0][0] == team) &&
		(g.board[1][1] == nil || *g.board[1][1] == team) &&
		(g.board[2][2] == nil || *g.board[2][2] == team) {
		score -= 10
	}
	if (g.board[0][2] == nil || *g.board[0][2] == team) &&
		(g.board[1][1] == nil || *g.board[1][1] == team) &&
		(g.board[2][0] == nil || *g.board[2][0] == team) {
		score -= 10
	}
	return score
}
