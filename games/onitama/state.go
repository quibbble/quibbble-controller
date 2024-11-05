package onitama

import (
	"fmt"
	"slices"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type state struct {
	seed    int64
	turn    string
	teams   []string
	winners []string
	board   [BoardSize][BoardSize]*Pawn
	hands   map[string][]string
	spare   string
}

func (s *state) Move(team string, startRow, startCol, endRow, endCol int, card string) error {
	if len(s.winners) > 0 {
		return fmt.Errorf("game already over")
	}
	if !slices.Contains(s.teams, team) {
		return fmt.Errorf("%s is not a valid team", team)
	}
	if team != s.turn {
		return fmt.Errorf("%s cannot play on %s turn", team, s.turn)
	}
	if !slices.Contains(s.hands[team], card) {
		return fmt.Errorf("%s hand does not contain %s", team, card)
	}
	if s.board[startRow][startCol] == nil {
		return fmt.Errorf("no pawn found at (%d, %d)", startRow, startCol)
	}
	if s.board[startRow][startCol] != nil && s.board[startRow][startCol].Team != team {
		return fmt.Errorf("cannot move enemy team's pawn")
	}
	if s.board[endRow][endCol] != nil && s.board[endRow][endCol].Team == team {
		return fmt.Errorf("cannot capture your own pawn")
	}
	if !canMove(card, endCol-startCol, endRow-startRow) {
		return fmt.Errorf("cannot move from (%d, %d) to {%d, %d} with %s", startRow, startCol, endRow, endCol, card)
	}
	// move the pawn
	captured := s.board[endRow][endCol]
	pawn := s.board[startRow][startCol]
	s.board[startRow][startCol] = nil
	s.board[endRow][endCol] = pawn

	// swap cards
	spare := s.spare
	s.spare = card
	s.hands[team][slices.Index(s.hands[team], card)] = spare

	// check for game over or change turn
	if (captured != nil && captured.Type == master) ||
		(team == s.teams[0] && endCol == 2 && endRow == BoardSize-1) ||
		(team == s.teams[1] && endCol == 2 && endRow == 0) {
		s.winners = []string{team}
	} else {
		s.turn = s.teams[(slices.Index(s.teams, team)+1)%len(s.teams)]
		// special case where player has no valid moves so turn is passed
		if len(s.actions()) == 0 {
			s.turn = s.teams[(slices.Index(s.teams, team)+1)%len(s.teams)]
		}
	}

	return nil
}

func canMove(card string, x, y int) bool {
	for _, movement := range movements[card] {
		if movement[0] == x && movement[1] == y {
			return true
		}
	}
	return false
}

func (s *state) actions() []*qg.Action {
	targets := make([]*qg.Action, 0)
	for startRow, row := range s.board {
		for startCol, pawn := range row {
			if pawn != nil && pawn.Team == s.turn {
				for _, card := range s.hands[s.turn] {
					for _, movement := range movements[card] {
						endRow := movement[1] + startRow
						endCol := movement[0] + startCol
						if (endRow >= 0 && endRow < BoardSize && endCol >= 0 && endCol < BoardSize) &&
							(s.board[endRow][endCol] == nil || s.board[endRow][endCol].Team != s.turn) {
							targets = append(targets, &qg.Action{
								Team: s.turn,
								Type: MoveAction,
								Details: MoveDetails{
									StartRow: startRow,
									StartCol: startCol,
									EndRow:   endRow,
									EndCol:   endCol,
									Card:     card,
								},
							})
						}
					}
				}
			}
		}
	}
	return targets
}

func (s *state) message() string {
	if len(s.winners) > 0 {
		return fmt.Sprintf("%s wins", s.winners[0])
	}
	return fmt.Sprintf("%s must move a pawn", s.turn)
}
