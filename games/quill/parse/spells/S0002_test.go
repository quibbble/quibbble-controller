package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0002(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0002")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 0, 0

	base := game.State.Board.XYs[x][y].Unit
	health := base.(*cd.UnitCard).Health

	if err := game.PlayCard(tests.Player1, uuids[0], base.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, health-1, game.State.Board.XYs[x][y].Unit.(*cd.UnitCard).Health)
}
