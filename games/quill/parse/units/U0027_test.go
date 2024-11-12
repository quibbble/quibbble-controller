package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0027(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0027")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	game.Board.XYs[x][y].Unit.(*cd.UnitCard).Movement = 1

	hand := game.Hand[tests.Player1].GetSize()

	if err := game.MoveUnit(tests.Player1, uuids[0], game.Board.XYs[x][y+1].UUID); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, hand+1, game.Hand[tests.Player1].GetSize())
}
