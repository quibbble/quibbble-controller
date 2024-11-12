package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_I0002(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "I0002")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1

	u2, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x][y+1].Unit = u2
	game.Board.XYs[x][y+1].Unit.(*cd.UnitCard).Cooldown = 0
	game.Board.XYs[x][y+1].Unit.(*cd.UnitCard).Codex = "11111111"

	health := u1.(*cd.UnitCard).Health

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, health+2, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health)
	assert.Equal(t, 1, len(u1.(*cd.UnitCard).GetTraits(tr.ShieldTrait)))

	if err := game.AttackUnit(tests.Player1, u2.GetUUID(), u1.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, health+1, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health)
}
