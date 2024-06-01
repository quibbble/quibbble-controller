package stratego

import (
	"testing"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

const (
	TeamA = "TeamA"
	TeamB = "TeamB"
)

func Test_Ai(t *testing.T) {
	stratego, err := NewStratego(ClassicVariant, 123, []string{TeamA, TeamB})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for i := 0; i < 10; i++ {
		err = stratego.Do(&qg.Action{
			Team: TeamA,
			Type: qg.AIAction,
		})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
