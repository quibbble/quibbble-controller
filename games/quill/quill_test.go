package quill

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Quill(t *testing.T) {
	game, err := NewQuill(123, []string{"A", "B"}, [][]string{
		{
			"S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001",
			"U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002",
			"U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002",
		},
		{
			"S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001",
			"U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002",
			"U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(game.teams) == 2)
	assert.True(t, len(game.targets) > 0)
}
