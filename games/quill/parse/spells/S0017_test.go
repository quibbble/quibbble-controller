package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0017(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0017")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1
	u2, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y+2].Unit = u2

	health := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health

	if err := game.PlayCard(tests.Player1, uuids[0]); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, health-1, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health)
	assert.Equal(t, health-1, game.Board.XYs[x][y+2].Unit.(*cd.UnitCard).Health)
}
