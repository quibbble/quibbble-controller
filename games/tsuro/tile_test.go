package tsuro

import (
	"testing"
)

func Test_NewTiel(t *testing.T) {
	edges := "CDEFGHAB"
	_, err := newTile(edges)
	if err != nil {
		t.FailNow()
	}
}
