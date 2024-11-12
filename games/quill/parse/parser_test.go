package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseCard(t *testing.T) {
	id := "U0001"
	card, err := ParseCard(id)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	assert.Equal(t, card.GetID(), id)
}
