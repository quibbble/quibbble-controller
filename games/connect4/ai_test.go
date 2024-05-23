package connect4

import (
	"testing"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	"github.com/stretchr/testify/assert"
)

func Test_Ai(t *testing.T) {
	connect4, err := NewConnect4([]string{TeamA, TeamB})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_ = connect4.Do(&qg.Action{Team: TeamA, Type: PlaceAction, Details: PlaceDetails{Col: 0}})
	_ = connect4.Do(&qg.Action{Team: TeamB, Type: PlaceAction, Details: PlaceDetails{Col: 0}})
	_ = connect4.Do(&qg.Action{Team: TeamA, Type: PlaceAction, Details: PlaceDetails{Col: 1}})
	_ = connect4.Do(&qg.Action{Team: TeamB, Type: PlaceAction, Details: PlaceDetails{Col: 1}})
	_ = connect4.Do(&qg.Action{Team: TeamA, Type: PlaceAction, Details: PlaceDetails{Col: 2}})
	_ = connect4.Do(&qg.Action{Team: TeamB, Type: PlaceAction, Details: PlaceDetails{Col: 2}})

	err = connect4.Do(&qg.Action{
		Team: TeamA,
		Type: qg.AIAction,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.True(t, len(connect4.winners) > 0)
}
