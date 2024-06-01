package tsuro

import (
	"testing"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	"github.com/stretchr/testify/assert"
)

const (
	TeamA = "TeamA"
	TeamB = "TeamB"
	TeamC = "TeamC"
	TeamD = "TeamD"
	TeamE = "TeamE"
	TeamF = "TeamF"
	TeamG = "TeamG"
	TeamH = "TeamH"
)

func Test_TsuroSmoke(t *testing.T) {
	tsuro, err := NewTsuro(ClassicVariant, time.Now().UnixNano(), []string{TeamA, TeamB})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	edges := "ABCDEFGH"
	rotated := "CDEFGHAB"
	tile, err := newTile(edges)
	if err != nil {
		t.FailNow()
	}
	tsuro.state.hands[TeamA].hand[0] = tile

	// rotate first tile in TeamA hand
	err = tsuro.Do(&qg.Action{
		Team: TeamA,
		Type: RotateAction,
		Details: RotateDetails{
			Tile: edges,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, rotated, tsuro.state.hands[TeamA].hand[0].Edges)

	tsuro.state.turn = TeamB
	tsuro.state.tokens[TeamB].Row = 0
	tsuro.state.tokens[TeamB].Col = 0
	tsuro.state.tokens[TeamB].Notch = "A"

	// place the first tile in TeamB hand at 0,0
	err = tsuro.Do(&qg.Action{
		Team: TeamB,
		Type: PlaceAction,
		Details: PlaceDetails{
			Row:  0,
			Col:  0,
			Tile: tsuro.state.hands[TeamB].hand[0].Edges,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
