package carcassonne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tile(t *testing.T) {
	t1 := newTile(City, Farm, City, City, NilStructure, false, true)

	t2 := newTile(City, Farm, City, City, NilStructure, false, true)
	assert.True(t, t1.equals(t2))

	t3 := newTile(Farm, City, City, City, NilStructure, false, true)
	assert.True(t, t1.equals(t3))

	t4 := newTile(City, City, City, Farm, NilStructure, false, true)
	assert.True(t, t1.equals(t4))
}
