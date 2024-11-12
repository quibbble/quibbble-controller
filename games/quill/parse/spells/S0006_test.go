package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0006(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0006")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1
	u2, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x+1][y].Unit = u2
	u3, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x-1][y].Unit = u3

	err = game.PlayCard(tests.Player1, uuids[0], u1.GetUUID(), u2.GetUUID())
	assert.True(t, err != nil)

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID(), u3.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, u3.GetUUID(), game.Board.XYs[x][y].Unit.GetUUID())
	assert.Equal(t, u1.GetUUID(), game.Board.XYs[x-1][y].Unit.GetUUID())
}
