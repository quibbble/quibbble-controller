package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0020(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0020")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}

	if err := game.EndTurn(tests.Player1); err != nil {
		t.Fatal(err)
	}

	codex := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Codex

	assert.True(t, codex != "00000000")

	if err := game.EndTurn(tests.Player2); err != nil {
		t.Fatal(err)
	}

	if err := game.EndTurn(tests.Player1); err != nil {
		t.Fatal(err)
	}

	assert.True(t, codex != game.Board.XYs[x][y].Unit.(*cd.UnitCard).Codex)
}
