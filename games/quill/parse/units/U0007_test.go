package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0007(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0007", "S0001", "S0010")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u1, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x][y-1].Unit = u1
	u2, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x][y+1].Unit = u2
	game.Board.XYs[x][y+1].Unit.(*cd.UnitCard).Cooldown = 0
	u3, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y+2].Unit = u3

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(game.Hooks()))

	cooldown := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown

	if err := game.PlayCard(tests.Player1, uuids[1], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, cooldown-1, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown)

	if err := game.AttackUnit(tests.Player1, u2.GetUUID(), u3.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(game.Hooks()))
	assert.Equal(t, cooldown-3, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown)

	if err := game.EndTurn(tests.Player1); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, cooldown-3, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown)

	if err := game.EndTurn(tests.Player2); err != nil {
		t.Fatal(err)
	}

	if err := game.PlayCard(tests.Player1, uuids[2], uuids[0]); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(game.Hooks()))
}
