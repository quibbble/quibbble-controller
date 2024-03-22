package stratego

import (
	"fmt"
	"math/rand"

	wr "github.com/mroth/weightedrand"
)

type Board struct {
	board [][]*Unit
}

func NewEmptyBoard(variant string) (*Board, error) {
	switch variant {
	case ClassicVariant:
		board := [][]*Unit{}
		for i := 0; i < BoardSize; i++ {
			row := []*Unit{}
			for j := 0; j < BoardSize; j++ {
				row = append(row, nil)
			}
			board = append(board, row)
		}
		for _, pair := range [][]int{{4, 2}, {4, 3}, {5, 2}, {5, 3}, {4, 6}, {4, 7}, {5, 6}, {5, 7}} {
			board[pair[0]][pair[1]] = Water()
		}
		return &Board{
			board: board,
		}, nil
	case QuickBattleVariant:
		board := [][]*Unit{}
		for i := 0; i < QuickBattleBoardSize; i++ {
			row := []*Unit{}
			for j := 0; j < QuickBattleBoardSize; j++ {
				row = append(row, nil)
			}
			board = append(board, row)
		}
		for _, pair := range [][]int{{3, 2}, {4, 2}, {3, 5}, {4, 5}} {
			board[pair[0]][pair[1]] = Water()
		}
		return &Board{
			board: board,
		}, nil
	default:
		return nil, fmt.Errorf("invalid variant %s", variant)
	}
}

func (b *Board) possibleMoves(row, col int) [][]int {
	boardSize := len(b.board)
	if row < 0 || row >= boardSize || col < 0 || col >= boardSize {
		return make([][]int, 0)
	}
	unit := b.board[row][col]
	if unit == nil || unit.Team == nil {
		return make([][]int, 0)
	}
	if unit.Type == bomb || unit.Type == flag {
		return make([][]int, 0)
	}
	moves := make([][]int, 0)
	if row+1 < boardSize && (b.board[row+1][col] == nil || (b.board[row+1][col].Team != nil && *b.board[row+1][col].Team != *unit.Team)) {
		moves = append(moves, []int{1, 0})
	}
	if row-1 >= 0 && (b.board[row-1][col] == nil || (b.board[row-1][col].Team != nil && *b.board[row-1][col].Team != *unit.Team)) {
		moves = append(moves, []int{-1, 0})
	}
	if col+1 < boardSize && (b.board[row][col+1] == nil || (b.board[row][col+1].Team != nil && *b.board[row][col+1].Team != *unit.Team)) {
		moves = append(moves, []int{0, 1})
	}
	if col-1 >= 0 && (b.board[row][col-1] == nil || (b.board[row][col-1].Team != nil && *b.board[row][col-1].Team != *unit.Team)) {
		moves = append(moves, []int{0, -1})
	}
	return moves
}

// retrieves the number of active (movable) units for a given team
func (b *Board) numActive(team string) int {
	count := 0
	for r, row := range b.board {
		for c, unit := range row {
			if unit != nil && unit.Team != nil &&
				*unit.Team == team &&
				len(b.possibleMoves(r, c)) > 0 {
				count += 1
			}
		}
	}
	return count
}

func NewRandomBoard(teams []string, variant string, random *rand.Rand) (*Board, error) {
	if len(teams) != 2 {
		return nil, fmt.Errorf("teams list must contain two teams")
	}
	unitToIdx := map[string]int{
		flag: 0, bomb: 1, spy: 2, scout: 3, miner: 4, sergeant: 5,
		lieutenant: 6, captain: 7, major: 8, colonel: 9, general: 10, marshal: 11,
	}
	idxToUnit := map[int]string{
		0: flag, 1: bomb, 2: spy, 3: scout, 4: miner, 5: sergeant,
		6: lieutenant, 7: captain, 8: major, 9: colonel, 10: general, 11: marshal,
	}

	var teamOneUnits, teamTwoUnits [12]int
	var flagChooser, bombChooser, minerChooser, scoutChooser *wr.Chooser

	switch variant {
	case ClassicVariant:
		teamOneUnits = [12]int{
			1, 6, 1, 8, 5, 4,
			4, 4, 3, 2, 1, 1,
		}
		teamTwoUnits = [12]int{
			1, 6, 1, 8, 5, 4,
			4, 4, 3, 2, 1, 1,
		}
		flagChooser, _ = wr.NewChooser(
			wr.Choice{Item: 2, Weight: 1},
			wr.Choice{Item: 1, Weight: 4},
			wr.Choice{Item: 0, Weight: 5},
		)
		bombChooser, _ = wr.NewChooser(
			wr.Choice{Item: 3, Weight: 1},
			wr.Choice{Item: 2, Weight: 2},
			wr.Choice{Item: 1, Weight: 3},
			wr.Choice{Item: 0, Weight: 4},
		)
		minerChooser, _ = wr.NewChooser(
			wr.Choice{Item: 2, Weight: 1},
			wr.Choice{Item: 1, Weight: 4},
			wr.Choice{Item: 0, Weight: 5},
		)
		scoutChooser, _ = wr.NewChooser(
			wr.Choice{Item: 3, Weight: 5},
			wr.Choice{Item: 2, Weight: 4},
			wr.Choice{Item: 1, Weight: 1},
		)
	case QuickBattleVariant:
		teamOneUnits = [12]int{
			1, 2, 1, 2, 2, 0,
			0, 0, 0, 0, 1, 1,
		}
		teamTwoUnits = [12]int{
			1, 2, 1, 2, 2, 0,
			0, 0, 0, 0, 1, 1,
		}
		flagChooser, _ = wr.NewChooser(
			wr.Choice{Item: 1, Weight: 5},
			wr.Choice{Item: 0, Weight: 5},
		)
		bombChooser, _ = wr.NewChooser(
			wr.Choice{Item: 2, Weight: 2},
			wr.Choice{Item: 1, Weight: 4},
			wr.Choice{Item: 0, Weight: 4},
		)
		minerChooser, _ = wr.NewChooser(
			wr.Choice{Item: 2, Weight: 1},
			wr.Choice{Item: 1, Weight: 4},
			wr.Choice{Item: 0, Weight: 5},
		)
		scoutChooser, _ = wr.NewChooser(
			wr.Choice{Item: 2, Weight: 5},
			wr.Choice{Item: 1, Weight: 4},
			wr.Choice{Item: 0, Weight: 1},
		)
	default:
		return nil, fmt.Errorf("invalid variant %s", variant)
	}

	board, err := NewEmptyBoard(variant)
	if err != nil {
		return nil, err
	}

	boardSize := len(board.board)

	place(board, flagChooser, random, true, NewUnit(flag, teams[0]))
	place(board, flagChooser, random, false, NewUnit(flag, teams[1]))
	teamOneUnits[unitToIdx[flag]] -= 1
	teamTwoUnits[unitToIdx[flag]] -= 1

	bombs := teamOneUnits[unitToIdx[bomb]]
	for i := 0; i < bombs; i++ {
		place(board, bombChooser, random, true, NewUnit(bomb, teams[0]))
		place(board, bombChooser, random, false, NewUnit(bomb, teams[1]))
		teamOneUnits[unitToIdx[bomb]] -= 1
		teamTwoUnits[unitToIdx[bomb]] -= 1
	}

	miners := teamOneUnits[unitToIdx[miner]]
	for i := 0; i < miners; i++ {
		place(board, minerChooser, random, true, NewUnit(miner, teams[0]))
		place(board, minerChooser, random, false, NewUnit(miner, teams[1]))
		teamOneUnits[unitToIdx[miner]] -= 1
		teamTwoUnits[unitToIdx[miner]] -= 1
	}

	scouts := teamOneUnits[unitToIdx[scout]]
	for i := 0; i < scouts; i++ {
		place(board, scoutChooser, random, true, NewUnit(scout, teams[0]))
		place(board, scoutChooser, random, false, NewUnit(scout, teams[1]))
		teamOneUnits[unitToIdx[scout]] -= 1
		teamTwoUnits[unitToIdx[scout]] -= 1
	}

	// place remainder randomly
	for row := 0; row < (boardSize-2)/2; row++ {
		for col := 0; col < boardSize; col++ {
			if board.board[row][col] == nil {
				for i, amt := range teamOneUnits {
					if amt > 0 {
						board.board[row][col] = NewUnit(idxToUnit[i], teams[0])
						teamOneUnits[i] -= 1
						break
					}
				}
			}
		}
	}
	for row := boardSize/2 + 1; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			if board.board[row][col] == nil {
				for i, amt := range teamTwoUnits {
					if amt > 0 {
						board.board[row][col] = NewUnit(idxToUnit[i], teams[1])
						teamTwoUnits[i] -= 1
						break
					}
				}
			}
		}
	}
	return board, nil
}

func getRandomNotTaken(board *Board, chooser *wr.Chooser, random *rand.Rand, isOne bool) (int, int) {
	boardSize := len(board.board)
	row := chooser.PickSource(random).(int)
	col := random.Intn(boardSize)
	if !isOne {
		row = boardSize - row - 1
	}
	if board.board[row][col] != nil {
		return getRandomNotTaken(board, chooser, random, isOne)
	}
	return row, col
}

func place(board *Board, chooser *wr.Chooser, random *rand.Rand, isOne bool, unit *Unit) {
	row, col := getRandomNotTaken(board, chooser, random, isOne)
	board.board[row][col] = unit
}
