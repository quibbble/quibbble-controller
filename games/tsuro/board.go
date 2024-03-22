package tsuro

import "fmt"

const (
	rows    = 6
	columns = 6
)

type board struct {
	board [][]*tile // 0,0 is top left corner
}

func newBoard() *board {
	var b = make([][]*tile, rows)
	for i := 0; i < rows; i++ {
		b[i] = make([]*tile, columns)
	}
	return &board{board: b}
}

func (b *board) Place(tile *tile, row, col int) error {
	if row < 0 || col < 0 || row >= rows || col >= columns {
		return fmt.Errorf("index out of bounds")
	}
	if b.board[row][col] != nil {
		return fmt.Errorf("tile already exists at (%d, %d)", row, col)
	}
	b.board[row][col] = tile
	return nil
}

func (b *board) getTileCount() int {
	counter := 0
	for _, row := range b.board {
		for _, tile := range row {
			if tile != nil {
				counter++
			}
		}
	}
	return counter
}
