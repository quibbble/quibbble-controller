package indigo

import (
	"testing"

	"github.com/quibbble/quibbble-controller/pkg/game"
)

const (
	TeamA = "TeamA"
	TeamB = "TeamB"
)

func Test_Indigo(t *testing.T) {
	indigo, err := NewIndigo(ClassicVariant, 123, []string{TeamA, TeamB})
	if err != nil {
		t.Fatal(err)
	}

	tile, _ := indigo.hands[TeamA].GetItem(0)

	if err := indigo.Do(&game.Action{
		Team: TeamA,
		Type: PlaceAction,
		Details: &PlaceDetails{
			Row:  4,
			Col:  5,
			Tile: tile.Paths,
		},
	}); err != nil {
		t.Fatal(err)
	}
}
