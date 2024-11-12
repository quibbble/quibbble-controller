package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0012(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0012", "S0001")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	mana := game.Mana[tests.Player1].BaseAmount

	if err := game.PlayCard(tests.Player1, uuids[1], uuids[0]); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, nil, game.Board.XYs[x][y].Unit)
	assert.Equal(t, mana+1, game.Mana[tests.Player1].BaseAmount)
}
