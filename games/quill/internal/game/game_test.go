package game

import (
	"testing"

	"github.com/quibbble/quibbble-controller/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_NewGame(t *testing.T) {
	p1 := uuid.UUID("P0000001")
	p2 := uuid.UUID("P0000002")

	deck1 := map[string]int{
		"U0002": 30,
	}
	deck2 := map[string]int{
		"U0002": 30,
	}

	d1 := make([]string, 0)
	d2 := make([]string, 0)

	for id, count := range deck1 {
		for i := 0; i < count; i++ {
			d1 = append(d1, id)
		}
	}
	for id, count := range deck2 {
		for i := 0; i < count; i++ {
			d2 = append(d2, id)
		}
	}

	game, err := NewGame(0, p1, p2, d1, d2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, game.GetTurn(), p1)
}
