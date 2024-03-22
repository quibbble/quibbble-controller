package stratego

import (
	"testing"
	"time"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	"github.com/stretchr/testify/assert"
)

func Test_StateRandomness(t *testing.T) {
	seed := time.Now().Unix()
	teams := []string{"A", "B"}

	s1, _ := newState(QuickBattleVariant, seed, teams)
	s2, _ := newState(QuickBattleVariant, seed, teams)

	for i := 0; i < len(s1.board.board); i++ {
		for j := 0; j < len(s1.board.board); j++ {
			u1 := s1.board.board[i][j]
			u2 := s2.board.board[i][j]
			if u1 != nil {
				assert.Equal(t, u1.Team, u2.Team)
				assert.Equal(t, u1.Type, u2.Type)
			}
		}
	}
}

func Test_BGN(t *testing.T) {
	raw := `
		[variant "quick_battle"]
		[key "stratego"]
		[teams "red, blue"]
		[seed "1696470849747"]`

	game, err := qgn.Parse(raw)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	builder := Builder{}
	if _, err := builder.Create(game); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
