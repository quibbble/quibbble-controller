package connect4

import (
	"fmt"
	"strings"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type state struct {
	turn    string
	teams   []string
	winners []string
	board   *board
}

func newState(teams []string) *state {
	return &state{
		turn:    teams[0],
		teams:   teams,
		winners: make([]string, 0),
		board:   newBoard(),
	}
}

func (s *state) Place(team string, col int) error {
	if len(s.winners) > 0 {
		return fmt.Errorf("game already over")
	}
	if team != s.turn {
		return fmt.Errorf("%s cannot play on %s turn", team, s.turn)
	}
	if col < 0 || col > Cols-1 {
		return fmt.Errorf("column %d is out of bounds", col)
	}
	if err := s.board.Place(team, col); err != nil {
		return err
	}
	if s.board.full() {
		s.winners = s.teams
		return nil
	}
	if winner := findWinner(s.board); winner != nil {
		s.winners = []string{*winner}
		return nil
	}
	for idx, team := range s.teams {
		if team == s.turn {
			s.turn = s.teams[(idx+1)%len(s.teams)]
			break
		}
	}
	return nil
}

// nil if no winner, winner name if winner
func findWinner(board *board) *string {
	// check vertical
	for row := 0; row < Rows-3; row++ {
		for col := 0; col < Cols; col++ {
			if board.board[row][col] == nil || board.board[row+1][col] == nil || board.board[row+2][col] == nil || board.board[row+3][col] == nil {
				continue
			}
			player := *board.board[row][col]
			if player == *board.board[row+1][col] && player == *board.board[row+2][col] && player == *board.board[row+3][col] {
				return &player
			}
		}
	}
	// check horizontal
	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols-3; col++ {
			if board.board[row][col] == nil || board.board[row][col+1] == nil || board.board[row][col+2] == nil || board.board[row][col+3] == nil {
				continue
			}
			player := *board.board[row][col]
			if player == *board.board[row][col+1] && player == *board.board[row][col+2] && player == *board.board[row][col+3] {
				return &player
			}
		}
	}
	// check positive diagonal
	for row := 0; row < Rows-3; row++ {
		for col := 0; col < Cols-3; col++ {
			if board.board[row][col] == nil || board.board[row+1][col+1] == nil || board.board[row+2][col+2] == nil || board.board[row+3][col+3] == nil {
				continue
			}
			player := *board.board[row][col]
			if player == *board.board[row+1][col+1] && player == *board.board[row+2][col+2] && player == *board.board[row+3][col+3] {
				return &player
			}
		}
	}
	// check negative diagonal
	for row := 3; row < Rows; row++ {
		for col := 0; col < Cols-3; col++ {
			if board.board[row][col] == nil || board.board[row-1][col+1] == nil || board.board[row-2][col+2] == nil || board.board[row-3][col+3] == nil {
				continue
			}
			player := *board.board[row][col]
			if player == *board.board[row-1][col+1] && player == *board.board[row-2][col+2] && player == *board.board[row-3][col+3] {
				return &player
			}
		}
	}
	return nil
}

func (s *state) actions() []*qg.Action {
	targets := make([]*qg.Action, 0)
	for col := 0; col < Cols; col++ {
		if s.board.board[0][col] == nil {
			targets = append(targets, &qg.Action{
				Team:    s.turn,
				Type:    PlaceAction,
				Details: PlaceDetails{Col: col},
			})
		}
	}
	return targets
}

func (s *state) message() string {
	message := fmt.Sprintf("%s must place a disk", s.turn)
	if len(s.winners) > 0 {
		message = fmt.Sprintf("%s tie", strings.Join(s.winners, " and "))
		if len(s.winners) == 1 {
			message = fmt.Sprintf("%s wins", s.winners[0])
		}
	}
	return message
}
