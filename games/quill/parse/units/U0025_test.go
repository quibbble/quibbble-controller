package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0025(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0025")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u1, _ := game.BuildCard("U0010", tests.Player1, false)
	game.Board.XYs[x][y+1].Unit = u1
	u2, _ := game.BuildCard("U0010", tests.Player2, false)
	game.Board.XYs[x+1][y].Unit = u2

	cooldown := game.Board.XYs[x][y+1].Unit.(*cd.UnitCard).Cooldown

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, cooldown-1, game.Board.XYs[x][y+1].Unit.(*cd.UnitCard).Cooldown)
	assert.Equal(t, cooldown-1, game.Board.XYs[x+1][y].Unit.(*cd.UnitCard).Cooldown)
}
