package state

import (
	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const (
	Rows = 7
	Cols = 5

	BaseID = "U0001"
)

type Tile struct {
	UUID uuid.UUID
	X, Y int
	Unit ICard
}

func NewTile(uuid uuid.UUID, x, y int) *Tile {
	return &Tile{
		UUID: uuid,
		X:    x,
		Y:    y,
	}
}

type Board struct {
	XYs   [Cols][Rows]*Tile
	UUIDs map[uuid.UUID]*Tile

	Sides map[uuid.UUID]int
}

func NewBoard(build BuildCard, gen *uuid.Gen, player1, player2 uuid.UUID) (*Board, error) {
	board := &Board{
		XYs:   [Cols][Rows]*Tile{},
		UUIDs: make(map[uuid.UUID]*Tile),
		Sides: make(map[uuid.UUID]int),
	}
	for x := 0; x < Cols; x++ {
		for y := 0; y < Rows; y++ {
			tile := NewTile(gen.New(en.TileUUID), x, y)
			board.XYs[x][y] = tile
			board.UUIDs[tile.UUID] = tile
		}
	}

	for x := 0; x < Cols; x++ {
		base1, err := build(BaseID, player1, false)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		base2, err := build(BaseID, player2, false)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		board.XYs[x][0].Unit = base1
		board.XYs[x][Rows-1].Unit = base2
	}

	board.Sides[player1] = 0
	board.Sides[player2] = Rows - 1

	return board, nil
}

func (b *Board) GetTileXY(tile uuid.UUID) (int, int, error) {
	for x, col := range b.XYs {
		for y, t := range col {
			if t.UUID == tile {
				return x, y, nil
			}
		}
	}
	return -1, -1, ErrNotFound(tile)
}

func (b *Board) GetUnitXY(unit uuid.UUID) (int, int, error) {
	for x, col := range b.XYs {
		for y, tile := range col {
			if tile.Unit != nil && tile.Unit.GetUUID() == unit {
				return x, y, nil
			}
		}
	}
	return -1, -1, ErrNotFound(unit)
}

// GetPlayableRowRange retries the range of rows the player may play a unit
func (b *Board) GetPlayableRowRange(player uuid.UUID) (int, int) {
	var min, max int
	if b.Sides[player] == 0 {
		min = 0
		max = 2

		full := true
	exit1:
		for x := 0; x < Cols; x++ {
			for y := min; y < max; y++ {
				if b.XYs[x][y].Unit == nil {
					full = false
					break exit1
				}
			}
		}
		if full {
			max++

			for x := 0; x < Cols; x++ {
				if b.XYs[x][max].Unit == nil {
					full = false
					break
				}
			}
			if full {
				max++
			}
		}
	} else {
		min = Rows - 3
		max = Rows - 1

		full := true
	exit2:
		for x := 0; x < Cols; x++ {
			for y := min; y < max; y++ {
				if b.XYs[x][y].Unit == nil {
					full = false
					break exit2
				}
			}
		}
		if full {
			min--

			for x := 0; x < Cols; x++ {
				if b.XYs[x][min].Unit == nil {
					full = false
					break
				}
			}
			if full {
				min--
			}
		}
	}
	return min, max
}
