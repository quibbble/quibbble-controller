package connect4

import (
	"math"
	"slices"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type AI struct{}

func (ai AI) Score(game qg.Game, team string) float64 {
	g := game.(*Connect4)
	if slices.Contains(g.winners, team) {
		return math.Inf(1)
	}
	if len(g.winners) > 0 {
		return math.Inf(-1)
	}
	score := 0.0

	open := func(token *string, team string) bool {
		return token == nil || *token == team
	}
	for i := range Rows - 4 {
		for j := range Cols - 4 {
			// down
			if open(g.board.board[i][j], team) &&
				open(g.board.board[i][j+1], team) &&
				open(g.board.board[i][j+2], team) &&
				open(g.board.board[i][j+3], team) {
				score -= 10
			}
			// right
			if open(g.board.board[i][j], team) &&
				open(g.board.board[i+1][j], team) &&
				open(g.board.board[i+2][j], team) &&
				open(g.board.board[i+3][j], team) {
				score -= 10
			}
			// diag
			if open(g.board.board[i][j], team) &&
				open(g.board.board[i+1][j+1], team) &&
				open(g.board.board[i+2][j+2], team) &&
				open(g.board.board[i+3][j+3], team) {
				score -= 10
			}
		}
	}
	return score
}
