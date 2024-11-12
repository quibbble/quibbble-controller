package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0019(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0019")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 5

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	assert.True(t, game.Board.XYs[x][y].Unit != nil)
	assert.Equal(t, 4, game.Board.XYs[x][y].Unit.GetCost())
}
