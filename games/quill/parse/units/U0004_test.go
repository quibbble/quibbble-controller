package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0004(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0004", "U0002", "S0003")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u0002, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x][y-1].Unit = u0002
	u0002Attack := u0002.(*cd.UnitCard).Attack

	// should play card and update friends traits
	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, game.Board.XYs[x][y].Unit.GetUUID(), uuids[0])
	assert.Equal(t, u0002Attack+1, game.Board.XYs[x][y-1].Unit.(*cd.UnitCard).Attack)

	// should play card and update friends traits
	if err := game.PlayCard(tests.Player1, uuids[1], game.Board.XYs[x+1][y].UUID); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, u0002Attack+1, game.Board.XYs[x+1][y].Unit.(*cd.UnitCard).Attack)

	// should move unit and update friends traits
	if err := game.MoveUnitXY(tests.Player1, uuids[0], x+1, y+1); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, u0002Attack, game.Board.XYs[x][y-1].Unit.(*cd.UnitCard).Attack)
	assert.Equal(t, u0002Attack+1, game.Board.XYs[x+1][y].Unit.(*cd.UnitCard).Attack)

	// should play card and update friends traits
	if err := game.PlayCard(tests.Player1, uuids[2], uuids[0]); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, u0002Attack, game.Board.XYs[x][y-1].Unit.(*cd.UnitCard).Health)
	assert.Equal(t, u0002Attack, game.Board.XYs[x+1][y].Unit.(*cd.UnitCard).Health)
}
