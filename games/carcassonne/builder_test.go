package carcassonne

import (
	"testing"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	"github.com/stretchr/testify/assert"
)

func Test_Builder_Notations(t *testing.T) {
	tags := map[string]string{
		"key":   Key,
		"teams": "TeamA, TeamB",
		"seed":  "123",
	}
	tests := []struct {
		name        string
		snapshot    *qgn.Snapshot
		shouldError bool
	}{
		{
			name:        "empty string should error",
			snapshot:    &qgn.Snapshot{},
			shouldError: true,
		},
		{
			name: "missing seed should error",
			snapshot: &qgn.Snapshot{
				Tags: map[string]string{
					"Game":  Key,
					"Teams": "TeamA, TeamB",
				},
			},
			shouldError: true,
		},
		{
			name: "should create a new game",
			snapshot: &qgn.Snapshot{
				Tags: tags,
			},
			shouldError: false,
		},
		{
			name: "should create a new game and do actions",
			snapshot: &qgn.Snapshot{
				Tags: tags,
				Actions: []qgn.Action{
					{
						Index:   0,
						Key:     "i",
						Details: []string{"1", "0", "f", "f", "r", "r", "n", "f", "f"},
					},
					{
						Index:   0,
						Key:     "o",
						Details: []string{"f", "1", "0", "f", "lb"},
					},
				},
			},
			shouldError: false,
		},
	}

	builder := Builder{}
	for _, test := range tests {
		_, err := builder.Create(test.snapshot)
		assert.Equal(t, test.shouldError, err != nil, test.name)
	}
}
