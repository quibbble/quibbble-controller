package units_tests

import (
	"testing"

	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0004(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0004")
	if err != nil {
		t.Fatal(err)
	}

	mana := game.Mana[tests.Player1].Amount
	handSize := game.Hand[tests.Player1].GetSize()

	if err := game.PlayCard(tests.Player1, uuids[0]); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, mana+3, game.Mana[tests.Player1].Amount)
	assert.Equal(t, handSize, game.Hand[tests.Player1].GetSize())
	assert.Equal(t, 0, len(game.Hooks()))
}
