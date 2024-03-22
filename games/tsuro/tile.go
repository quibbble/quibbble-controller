package tsuro

import (
	"fmt"
	"slices"
)

/*
tile representation

	     A  B
	    ——  ——
	H |        | C
	G |        | D
	    ——  ——
	     F  E
*/
type tile struct {
	Edges string            `json:"edges"` // defines the tile
	Paths map[string]string `json:"paths"` // map from path to team defining section team path
}

func newTile(edges string) (*tile, error) {
	t := &tile{
		Edges: edges,
		Paths: make(map[string]string),
	}
	for i := 0; i < 4; i++ {
		if slices.Contains(tiles, t.Edges) {
			return &tile{
				Edges: edges,
				Paths: make(map[string]string),
			}, nil
		}
		t.RotateRight()
	}
	return nil, fmt.Errorf("edges %s are not a valid tile configuration", edges)
}

func (t *tile) GetDestination(start string) string {
	for idx, char := range t.Edges {
		if string(char) == start && idx%2 == 0 {
			return string(t.Edges[idx+1])
		} else if string(char) == start && idx%2 == 1 {
			return string(t.Edges[idx-1])
		}
	}
	return ""
}

func (t *tile) RotateRight() {
	transform := map[string]string{"A": "C", "B": "D", "C": "E", "D": "F", "E": "G", "F": "H", "G": "A", "H": "B"}
	transformed := ""
	for _, char := range t.Edges {
		transformed += transform[string(char)]
	}
	t.Edges = transformed
}

func (t *tile) RotateLeft() {
	transform := map[string]string{"A": "G", "B": "H", "C": "A", "D": "B", "E": "C", "F": "D", "G": "E", "H": "F"}
	transformed := ""
	for _, char := range t.Edges {
		transformed += transform[string(char)]
	}
	t.Edges = transformed
}

func (t *tile) countCrossings(team string) int {
	paths := make([]string, 0)
	for path, t := range t.Paths {
		if team == t {
			paths = append(paths, path)
		}
	}
	count := 0
	for i, path := range paths {
		for j, otherPath := range paths {
			if i != j && slices.Contains(crossing[path], otherPath) {
				count++
			}
		}
	}
	return count / 2
}

func (t *tile) equals(t2 *tile) bool {
	copied, _ := newTile(t.Edges)
	for i := 0; i < 4; i++ {
		if copied.Edges == t2.Edges {
			return true
		}
		copied.RotateRight()
	}
	return false
}

func (t *tile) in(list []*tile) bool {
	for _, t2 := range list {
		if t.equals(t2) {
			return true
		}
	}
	return false
}
