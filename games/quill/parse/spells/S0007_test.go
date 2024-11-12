package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0007(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0007")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1
	handSize := game.Hand[tests.Player1].GetSize()

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	newX, newY, err := game.Board.GetUnitXY(u1.GetUUID())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, handSize, game.Hand[tests.Player1].GetSize())
	assert.True(t, x != newX || y != newY)
}
