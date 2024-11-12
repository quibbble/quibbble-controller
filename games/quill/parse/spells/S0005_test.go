package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0005(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0005")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1
	u2, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x+1][y].Unit = u2
	u3, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x-1][y].Unit = u3
	u4, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x-1][y+1].Unit = u4
	u5, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y+4].Unit = u5

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.True(t, game.Board.XYs[x][y].Unit == nil)
	assert.True(t, game.Board.XYs[x+1][y].Unit == nil)
	assert.True(t, game.Board.XYs[x-1][y].Unit == nil)
	assert.True(t, game.Board.XYs[x-1][y+1].Unit == nil)
	assert.True(t, game.Board.XYs[x][y+4].Unit != nil)
}
