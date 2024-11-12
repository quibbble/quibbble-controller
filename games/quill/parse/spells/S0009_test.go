package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0009(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0009")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1
	handSize := game.Hand[tests.Player1].GetSize()

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1

	err = game.PlayCard(tests.Player1, uuids[0], u1.GetUUID())
	assert.True(t, err != nil)

	game.Board.XYs[x][y].Unit.(*cd.UnitCard).Player = tests.Player1

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, handSize+2, game.Hand[tests.Player1].GetSize())
	assert.True(t, game.Board.XYs[x][y].Unit == nil)
}
