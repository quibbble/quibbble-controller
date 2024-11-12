package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0022(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0022")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u1, _ := game.BuildCard("U0010", tests.Player2, false)
	game.Board.XYs[x][y+2].Unit = u1
	u2, _ := game.BuildCard("U0010", tests.Player2, false)
	game.Board.XYs[x][y+3].Unit = u2

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown = 0

	if err := game.AttackUnit(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, nil, game.Board.XYs[x][y+2].Unit)
	assert.Equal(t, nil, game.Board.XYs[x][y+3].Unit)
}
