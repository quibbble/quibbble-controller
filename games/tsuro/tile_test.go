package tsuro

import (
	"testing"
)

func Test_NewTile(t *testing.T) {
	edges := "CDEFGHAB"
	_, err := newTile(edges)
	if err != nil {
		t.FailNow()
	}
}
