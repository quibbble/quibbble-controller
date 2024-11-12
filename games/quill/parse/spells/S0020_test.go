package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0020(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0020")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1

	cooldown := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, cooldown-1, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown)
}
