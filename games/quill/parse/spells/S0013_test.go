package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0013(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0013")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1

	startCount := 0
	for _, col := range game.Board.XYs {
		for _, tile := range col {
			if tile.Unit != nil && tile.Unit.GetUUID() != u1.GetUUID() {
				startCount++
			}
		}
	}

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	count := 0
	for _, col := range game.Board.XYs {
		for _, tile := range col {
			if tile.Unit != nil && tile.Unit.GetUUID() != u1.GetUUID() {
				count++
			}
		}
	}
	assert.Equal(t, startCount+1, count)
}
