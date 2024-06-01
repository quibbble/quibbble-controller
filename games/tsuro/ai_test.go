package tsuro

import (
	"testing"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

func Test_Ai(t *testing.T) {
	tsuro, err := NewTsuro(ClassicVariant, 4321, []string{TeamA, TeamB, TeamC, TeamD, TeamE, TeamF, TeamG, TeamH})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for i := 0; i < 20; i++ {
		err = tsuro.Do(&qg.Action{
			Team: TeamA,
			Type: qg.AIAction,
		})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
