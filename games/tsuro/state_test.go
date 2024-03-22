package tsuro

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewState(t *testing.T) {
	testCases := []struct {
		name      string
		teams     []string
		seed      int64
		variant   string
		shouldErr bool
	}{
		{
			name:      "teams greater than 11 should error",
			teams:     []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},
			seed:      123,
			variant:   ClassicVariant,
			shouldErr: true,
		},
		{
			name:      "duplicate teams should error",
			teams:     []string{"1", "2", "1"},
			seed:      123,
			variant:   ClassicVariant,
			shouldErr: true,
		},
		{
			name:      "invalid variant should error",
			teams:     []string{"1", "2"},
			seed:      123,
			variant:   "VariantInvalid",
			shouldErr: true,
		},
		{
			name:      "missing variant should error",
			teams:     []string{"1", "2"},
			seed:      123,
			variant:   "",
			shouldErr: true,
		},
	}
	for _, test := range testCases {
		_, err := newState(test.variant, test.seed, test.teams)
		assert.Equal(t, err != nil, test.shouldErr, "ERROR: ", test.name)
	}
}
