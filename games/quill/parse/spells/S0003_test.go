package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0003(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0003")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u0002, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u0002

	handSize := game.Hand[tests.Player1].GetSize()

	if err := game.PlayCard(tests.Player1, uuids[0], u0002.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.True(t, game.Board.XYs[x][y].Unit == nil)
	assert.Equal(t, handSize, game.Hand[tests.Player1].GetSize())
	assert.Equal(t, 0, len(game.Hooks()))
}
