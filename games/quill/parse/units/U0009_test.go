package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0009(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0009")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(game.Hooks()))

	handSize := game.Hand[tests.Player1].GetSize()

	turns := 3
	for i := 0; i < turns; i++ {

		assert.True(t, game.Board.XYs[x][y].Unit != nil)

		if err := game.EndTurn(tests.Player1); err != nil {
			t.Fatal(err)
		}
		if err := game.EndTurn(tests.Player2); err != nil {
			t.Fatal(err)
		}
	}
	assert.Equal(t, nil, game.Board.XYs[x][y].Unit)
	assert.Equal(t, handSize+turns+2, game.Hand[tests.Player1].GetSize())
	assert.Equal(t, 0, len(game.Hooks()))
}
