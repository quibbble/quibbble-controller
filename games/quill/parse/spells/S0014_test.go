package units_tests

import (
	"testing"

	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0014(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0014", "I0004")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1

	if err := game.PlayCard(tests.Player1, uuids[1], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(game.Board.XYs[x][y].Unit.GetTraits(tr.PoisonTrait)))

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, len(game.Board.XYs[x][y].Unit.GetTraits(tr.PoisonTrait)))
}
