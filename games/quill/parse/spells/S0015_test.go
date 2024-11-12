package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0015(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0015")
	if err != nil {
		t.Fatal(err)
	}

	hand := game.Hand[tests.Player1].GetSize()

	if err := game.PlayCard(tests.Player1, uuids[0]); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, hand+1, game.Hand[tests.Player1].GetSize())
}
