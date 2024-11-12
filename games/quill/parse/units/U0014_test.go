package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0014(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0014", "I0001")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y-1].Unit = u1

	if err := game.PlayCard(tests.Player1, uuids[1], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown = 0

	attack := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Attack

	if err := game.AttackUnit(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(game.Board.XYs[x][y].Unit.(*cd.UnitCard).Items))
	assert.Equal(t, 0, len(game.Board.XYs[x][y-1].Unit.(*cd.UnitCard).Items))
	assert.True(t, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Attack > attack)
}
