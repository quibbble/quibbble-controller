package tictactoe

import (
	"fmt"
	"slices"
	"strings"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type state struct {
	turn    string
	teams   []string
	winners []string
	board   [BoardSize][BoardSize]*string
}

func (s *state) Mark(team string, row, col int) error {
	if len(s.winners) > 0 {
		return fmt.Errorf("game has ended")
	}
	if !slices.Contains(s.teams, team) {
		return fmt.Errorf("%s is not a valid team", team)
	}
	if team != s.turn {
		return fmt.Errorf("%s cannot play on %s turn", team, s.turn)
	}
	if row < 0 || row >= BoardSize || col < 0 || col >= BoardSize {
		return fmt.Errorf("row or column out of bounds")
	}
	if s.board[row][col] != nil {
		return fmt.Errorf("(%d,%d) already marked", row, col)
	}

	s.board[row][col] = &team

	if winner := checkWinner(s.board); winner != nil {
		s.winners = []string{*winner}
	} else if checkDraw(s.board) {
		s.winners = s.teams
	}

	s.turn = s.teams[(slices.Index(s.teams, team)+1)%len(s.teams)]
	return nil
}

func checkWinner(board [BoardSize][BoardSize]*string) *string {
	equal := func(a, b *string) bool {
		if a != nil && b != nil && *a == *b {
			return true
		}
		return false
	}

	for i := 0; i < BoardSize; i++ {
		if equal(board[i][0], board[i][1]) && equal(board[i][0], board[i][2]) {
			return board[i][0]
		}
		if equal(board[0][i], board[1][i]) && equal(board[0][i], board[2][i]) {
			return board[0][i]
		}
	}
	if equal(board[0][0], board[1][1]) && equal(board[0][0], board[2][2]) {
		return board[0][0]
	}
	if equal(board[2][0], board[1][1]) && equal(board[2][0], board[0][2]) {
		return board[2][0]
	}
	return nil
}

func checkDraw(board [BoardSize][BoardSize]*string) bool {
	for _, row := range board {
		for _, loc := range row {
			if loc == nil {
				return false
			}
		}
	}
	return true
}

func (s *state) getActions() []*qg.Action {
	targets := make([]*qg.Action, 0)
	for r, row := range s.board {
		for c, loc := range row {
			if loc == nil {
				targets = append(targets, &qg.Action{
					Team: s.turn,
					Type: MarkAction,
					Details: MarkDetails{
						Row: r,
						Col: c,
					},
				})
			}
		}
	}
	return targets
}

func (s *state) getMessage() string {
	message := fmt.Sprintf("%s must mark a location", s.turn)
	if len(s.winners) > 0 {
		message = fmt.Sprintf("%s tie", strings.Join(s.winners, " and "))
		if len(s.winners) == 1 {
			message = fmt.Sprintf("%s wins", s.winners[0])
		}
	}
	return message
}
