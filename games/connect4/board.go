package connect4

import "fmt"

type board struct {
	board [Rows][Cols]*string // 0,0 is top left corner
}

func newBoard() *board {
	return &board{
		board: [Rows][Cols]*string{},
	}
}

func (b *board) Place(team string, col int) error {
	for row := len(b.board) - 1; row >= 0; row-- {
		if b.board[row][col] == nil {
			b.board[row][col] = &team
			return nil
		}
	}
	return fmt.Errorf("column %d is full", col)
}

func (b *board) full() bool {
	for _, row := range b.board {
		for _, it := range row {
			if it == nil {
				return false
			}
		}
	}
	return true
}
