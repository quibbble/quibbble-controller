package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0011(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0011")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u1, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x][y+1].Unit = u1
	game.Board.XYs[x][y+1].Unit.(*cd.UnitCard).Attack = 1

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown = 0

	if err := game.AttackUnit(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.True(t, game.Board.XYs[x][y].Unit == nil)

	_, _, err = game.Board.GetUnitXY(uuids[0])
	assert.True(t, err == nil)
}
