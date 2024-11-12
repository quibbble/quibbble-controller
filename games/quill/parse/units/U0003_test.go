package units_tests

import (
	"testing"

	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	"github.com/quibbble/quibbble-controller/games/quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0003(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0003")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 0, 1

	// should play card
	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, game.Board.XYs[x][y].Unit.GetUUID(), uuids[0])

	u0002, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y+2].Unit = u0002

	// should fail cooldown check
	err = game.AttackUnit(tests.Player1, uuids[0], u0002.GetUUID())
	assert.True(t, err != nil)

	targets, err := game.GetNextTargets(tests.Player1, uuids[0])
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(targets))

	// should attack at range
	u0002Health := game.Board.XYs[x][y+2].Unit.(*cd.UnitCard).Health
	u0003Health := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health
	game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown = 0

	targets, err = game.GetNextTargets(tests.Player1, uuids[0])
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(targets))

	if err := game.MoveUnit(tests.Player1, uuids[0], game.Board.XYs[x+1][y].UUID); err != nil {
		t.Fatal(err)
	}

	targets, err = game.GetNextTargets(tests.Player1, uuids[0])
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(targets))

	targets, err = game.GetNextTargets(tests.Player1, uuids[0], u0002.GetUUID())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(targets))

	if err := game.AttackUnit(tests.Player1, uuids[0], u0002.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, u0003Health, game.Board.XYs[x+1][y].Unit.(*cd.UnitCard).Health)
	assert.Equal(t, u0002Health-1, game.Board.XYs[x][y+2].Unit.(*cd.UnitCard).Health)
}
