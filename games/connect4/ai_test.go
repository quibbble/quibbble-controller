package connect4

import (
	"testing"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

func Test_Ai(t *testing.T) {
	connect4, err := NewConnect4([]string{TeamA, TeamB})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	connect4.state.turn = TeamA

	err = connect4.Do(&qg.Action{
		Team: TeamA,
		Type: qg.AIAction,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
