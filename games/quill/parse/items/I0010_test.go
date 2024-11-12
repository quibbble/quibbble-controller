package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_I0010(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "I0010")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x][y].Unit = u1

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "11111111", game.Board.XYs[x][y].Unit.(*cd.UnitCard).Codex)
}
