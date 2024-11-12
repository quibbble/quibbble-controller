package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0026(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0026", "S0001")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	attack := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Attack

	if err := game.PlayCard(tests.Player1, uuids[1], uuids[0]); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, attack+1, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Attack)
}
