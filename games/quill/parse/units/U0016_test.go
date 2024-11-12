package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
)

func Test_U0016(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0016")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}
}
