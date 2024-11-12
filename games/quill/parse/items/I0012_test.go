package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_I0012(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "I0012", "S0001")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1

	health := u1.(*cd.UnitCard).Health

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(u1.(*cd.UnitCard).GetTraits(tr.WardTrait)))

	if err := game.PlayCard(tests.Player1, uuids[1], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, health, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health)
}
