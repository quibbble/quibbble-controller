package connect4

import (
	"testing"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	"github.com/stretchr/testify/assert"
)

const (
	TeamA = "TeamA"
	TeamB = "TeamB"
)

func Test_Connect4(t *testing.T) {
	connect4, err := NewConnect4([]string{TeamA, TeamB})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	connect4.state.turn = TeamA

	// place disk in column 0
	err = connect4.Do(&qg.Action{
		Team: TeamA,
		Type: PlaceAction,
		Details: PlaceDetails{
			Col: 0,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, TeamA, *connect4.state.board.board[Rows-1][0])
}
