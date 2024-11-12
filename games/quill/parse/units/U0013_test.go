package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0013(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0013", "S0001")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y-1].Unit = u1

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, game.Mana[tests.Player1].Amount, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Attack)

	if err := game.PlayCard(tests.Player1, uuids[1], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, game.Mana[tests.Player1].Amount, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Attack)
}
