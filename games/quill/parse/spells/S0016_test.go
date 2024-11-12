package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/quibbble/quibbble-controller/pkg/maths"
	"github.com/stretchr/testify/assert"
)

func Test_S0016(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0016")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1

	movement := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Movement

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, maths.MaxInt(0, movement+2), game.Board.XYs[x][y].Unit.(*cd.UnitCard).Movement)
}
