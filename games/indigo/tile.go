package indigo

import (
	"fmt"
	"slices"
)

// Edges
const (
	A       = "A"
	B       = "B"
	C       = "C"
	D       = "D"
	E       = "E"
	F       = "F"
	Special = "S" // Special edge represents all edges on the central treasure tile
)

/*
tile representation

	     A  B
	     /  \
	  F |    | C
		 \  /
		 E  D
*/
type tile struct {
	Paths    string // defines the paths on the tile - EX: ABCDEF means A <> B, C <> D, E <> F
	Treasure bool   // spsecial tiles with different paths rules
}

func newTreasureTile(paths string) *tile {
	return &tile{
		Paths:    paths,
		Treasure: true,
	}
}

func newTile(paths string) (*tile, error) {
	t := &tile{
		Paths: paths,
	}
	for i := 0; i < 6; i++ {
		if slices.Contains(uniquePaths, t.Paths) {
			return &tile{
				Paths: paths,
			}, nil
		}
		t.RotateClockwise()
	}
	return nil, fmt.Errorf("paths %s are not a valid tile configuration", paths)
}

func (t *tile) GetDestination(startingEdge string) (string, error) {
	for idx, char := range t.Paths {
		if string(char) == startingEdge && idx%2 == 0 {
			return string(t.Paths[idx+1]), nil
		} else if string(char) == startingEdge && idx%2 == 1 {
			return string(t.Paths[idx-1]), nil
		}
	}
	return "", fmt.Errorf("no destination found for tile %s with starting edge %s", t.Paths, startingEdge)
}

func (t *tile) RotateClockwise() {
	transform := map[string]string{A: B, B: C, C: D, D: E, E: F, F: A}
	transformed := ""
	for _, char := range t.Paths {
		transformed += transform[string(char)]
	}
	t.Paths = transformed
}

func (t *tile) equals(t2 *tile) bool {
	copied, _ := newTile(t.Paths)
	for i := 0; i < 6; i++ {
		if copied.Paths == t2.Paths {
			return true
		}
		copied.RotateClockwise()
	}
	return false
}
