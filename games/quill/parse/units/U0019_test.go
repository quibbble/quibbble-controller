package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0019(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0019")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u1, _ := game.BuildCard("U0010", tests.Player2, false)
	game.Board.XYs[x][y-1].Unit = u1

	attack := u1.(*cd.UnitCard).Attack

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, attack-1, game.Board.XYs[x][y-1].Unit.(*cd.UnitCard).Attack)
}
