package carcassonne

import (
	"testing"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

func Test_Ai(t *testing.T) {
	carcassonne, err := NewCarcassonne(123, []string{TeamA, TeamB})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = carcassonne.Do(&qg.Action{
		Team: TeamA,
		Type: qg.AIAction,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
