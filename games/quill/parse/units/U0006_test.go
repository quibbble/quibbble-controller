package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0006(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0006", "S0001")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u0006, _ := game.Hand[tests.Player1].GetCard(uuids[0])
	mana := game.Mana[tests.Player1].Amount - u0006.GetCost()
	baseMana := game.Mana[tests.Player1].BaseAmount

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, mana+1, game.Mana[tests.Player1].Amount)
	assert.Equal(t, baseMana+1, game.Mana[tests.Player1].BaseAmount)

	s0001, _ := game.Hand[tests.Player1].GetCard(uuids[1])
	mana = game.Mana[tests.Player1].Amount - s0001.GetCost()
	baseMana = game.Mana[tests.Player1].BaseAmount

	if err := game.PlayCard(tests.Player1, uuids[1], game.Board.XYs[x][y].Unit.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, nil, game.Board.XYs[x][y].Unit)
	assert.Equal(t, mana-1, game.Mana[tests.Player1].Amount)
	assert.Equal(t, baseMana-1, game.Mana[tests.Player1].BaseAmount)
}
