package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0011(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0011")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0005", tests.Player1, false)
	game.Board.XYs[x][y].Unit = u1

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, nil, game.Board.XYs[x][y].Unit)
	assert.Equal(t, "U0005", game.Hand[tests.Player1].GetItems()[game.Hand[tests.Player1].GetSize()-1].GetID())
}
