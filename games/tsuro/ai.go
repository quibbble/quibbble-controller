package tsuro

import (
	"math"
	"slices"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type AI struct{}

func (ai AI) Score(game qg.Game, team string) float64 {
	g := game.(*Tsuro)
	if slices.Contains(g.winners, team) {
		return math.Inf(1)
	}
	if len(g.winners) > 0 {
		return math.Inf(-1)
	}
	score := 0.0

	// find number of open spaces connected to your tile and give +10 for each space
	t, _ := g.tokens[g.turn].getAdjacent()
	row := t.Row
	col := t.Col

	type node struct {
		row, col int
	}
	s, f := []node{}, []node{}
	s = append(s, node{row, col})
	for len(s) > 0 {
		v := s[len(s)-1]
		s = s[:len(s)-1]
		if !slices.Contains(f, v) {
			f = append(f, v)
			if v.row+1 < BoardSize && g.board.board[v.row+1][v.col] == nil {
				s = append(s, node{v.row + 1, v.col})
			}
			if v.row-1 >= 0 && g.board.board[v.row-1][v.col] == nil {
				s = append(s, node{v.row - 1, v.col})
			}
			if v.col+1 < BoardSize && g.board.board[v.row][v.col+1] == nil {
				s = append(s, node{v.row, v.col + 1})
			}
			if v.col-1 >= 0 && g.board.board[v.row][v.col-1] == nil {
				s = append(s, node{v.row, v.col - 1})
			}
		}
	}

	score += float64(len(f) * 10)
	return score
}
