package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0012(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0012")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1
	game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown = 0
	u2, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x][y+1].Unit = u2

	cooldown1 := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown
	cooldown2 := game.Board.XYs[x][y+1].Unit.(*cd.UnitCard).Cooldown

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID(), u2.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, cooldown2, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown)
	assert.Equal(t, cooldown1, game.Board.XYs[x][y+1].Unit.(*cd.UnitCard).Cooldown)
}
