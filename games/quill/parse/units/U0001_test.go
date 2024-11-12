package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0001(t *testing.T) {
	game, _, err := tests.NewTestEnv(tests.Player1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, game.Board.XYs[0][0].Unit.GetID(), "U0001")
}
