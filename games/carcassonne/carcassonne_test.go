package carcassonne

import (
	"testing"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	"github.com/stretchr/testify/assert"
)

const (
	TeamA = "TeamA"
	TeamB = "TeamB"
)

func Test_Carcassonne(t *testing.T) {
	carcassonne, err := NewCarcassonne(time.Now().UnixNano(), []string{TeamA, TeamB})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 1, len(carcassonne.state.board.board), "board missing start tile")

	// place all road tile to right of start tile
	carcassonne.state.playTiles[carcassonne.state.turn] = newTile(Road, Road, Road, Road, NilStructure, false, false)
	carcassonne.state.turn = TeamA
	err = carcassonne.Do(&qg.Action{
		Team: TeamA,
		Type: PlaceTileAction,
		Details: PlaceTileDetails{
			X: 1,
			Y: 0,
			Tile: Tile{
				Road, Road, Road, Road, NilStructure, false, false,
			},
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 2, len(carcassonne.state.board.board), "board missing last placed tile")

	// place token in left B side of farmland
	err = carcassonne.Do(&qg.Action{
		Team: TeamA,
		Type: PlaceTokenAction,
		Details: PlaceTokenDetails{
			Pass: false,
			X:    1,
			Y:    0,
			Type: Farmer,
			Side: FarmSideLeftB,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 1, len(carcassonne.state.boardTokens), "missing placed token")
	assert.Equal(t, Farmer, carcassonne.state.boardTokens[0].Type, "incorrect token placed")
	assert.Equal(t, TeamB, carcassonne.state.turn, "incorrect team's turn")

	// place tile to left of start tile completing a road segment
	carcassonne.state.playTiles[carcassonne.state.turn] = newTile(Road, Road, Farm, Road, NilStructure, false, false)
	err = carcassonne.Do(&qg.Action{
		Team: TeamB,
		Type: PlaceTileAction,
		Details: PlaceTileDetails{
			X: -1,
			Y: 0,
			Tile: Tile{
				Road, Road, Farm, Road, NilStructure, false, false,
			},
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 3, len(carcassonne.state.board.board), "board missing last place tile")

	// claim the completed road by placing thief on right side of tile
	err = carcassonne.Do(&qg.Action{
		Team: TeamB,
		Type: PlaceTokenAction,
		Details: PlaceTokenDetails{
			Pass: false,
			X:    -1,
			Y:    0,
			Type: Thief,
			Side: SideRight,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 1, len(carcassonne.state.boardTokens))
	assert.Equal(t, 3, carcassonne.state.scores[TeamB])
	assert.Equal(t, []string{TeamB}, carcassonne.state.board.board[0].Teams[SideLeft])
	assert.Equal(t, []string{TeamB}, carcassonne.state.board.board[0].Teams[SideRight])
	assert.Equal(t, []string{TeamB}, carcassonne.state.board.board[1].Teams[SideLeft])
	assert.Equal(t, []string{TeamB}, carcassonne.state.board.board[2].Teams[SideRight])
	assert.Equal(t, TeamA, carcassonne.state.turn)
}
