package gamenotation

import (
	"fmt"
	"strconv"
	"strings"
)

type Tags map[string]string

const (
	KeyTag     = "key"     // unique name of the game being played
	IDTag      = "id"      // unique game id
	TeamsTag   = "teams"   // list of teams playing the game
	VariantTag = "variant" // game variant
	SeedTag    = "seed"    // random seed used to generate deterministic randomness
)

var RequiredTags = []string{KeyTag, TeamsTag}

func (t Tags) Teams() ([]string, error) {
	raw, ok := t[TeamsTag]
	if !ok {
		return nil, fmt.Errorf("required %s tag is missing", TeamsTag)
	}
	teams := strings.Split(raw, ", ")
	return teams, nil
}

func (t Tags) Seed() (int64, error) {
	raw, ok := t[SeedTag]
	if !ok {
		return -1, fmt.Errorf("optional %s tag is missing", SeedTag)
	}
	i, err := strconv.Atoi(raw)
	return int64(i), err
}
