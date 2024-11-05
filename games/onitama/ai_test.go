package onitama

import (
	"testing"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

func Test_AI(t *testing.T) {
	g, err := NewOnitama(123, []string{"red", "blue"})
	if err != nil {
		t.Fatal(err)
	}
	if err := qg.AI(Builder{}, AI{}, g, 3); err != nil {
		t.Fatal(err)
	}
}
