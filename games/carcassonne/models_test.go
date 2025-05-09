package carcassonne

import (
	"encoding/json"
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Mapstructure_Decoding(t *testing.T) {
	raw := `{"connected_cities": true}`

	a := map[string]interface{}{}

	if err := json.Unmarshal([]byte(raw), &a); err != nil {
		t.Error(err)
		t.FailNow()
	}

	var tile Tile
	if err := mapstructure.Decode(a, &tile); err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.True(t, tile.ConnectedCities)
}
