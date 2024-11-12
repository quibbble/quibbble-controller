package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_I0008(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "I0008", "S0001", "S0001")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x][y].Unit = u1
	u2, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x][y+1].Unit = u2
	u3, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y+2].Unit = u3

	attack := u1.(*cd.UnitCard).Attack

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(game.Board.XYs[x][y].Unit.(*cd.UnitCard).Items))
	assert.Equal(t, 1, len(game.Board.XYs[x][y].Unit.(*cd.UnitCard).GetTraits(tr.BuffTrait)))
	assert.Equal(t, attack+3, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Attack)

	if err := game.PlayCard(tests.Player1, uuids[1], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, nil, game.Board.XYs[x][y].Unit)
	assert.Equal(t, 1, len(game.Board.XYs[x][y+1].Unit.(*cd.UnitCard).Items))
	assert.Equal(t, 1, len(game.Board.XYs[x][y+1].Unit.(*cd.UnitCard).GetTraits(tr.BuffTrait)))

	if err := game.PlayCard(tests.Player1, uuids[2], u2.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, nil, game.Board.XYs[x][y+1].Unit)
	assert.Equal(t, 0, len(game.Board.XYs[x][y+2].Unit.(*cd.UnitCard).Items))
	assert.Equal(t, 0, len(game.Board.XYs[x][y+2].Unit.(*cd.UnitCard).GetTraits(tr.BuffTrait)))
}
