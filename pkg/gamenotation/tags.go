package gamenotation

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Tags map[string]string

const (
	KeyTag     = "key"     // unique name of the game being played
	TeamsTag   = "teams"   // list of teams playing the game
	VariantTag = "variant" // game variant
	SeedTag    = "seed"    // random seed used to generate deterministic randomness

	// server only tags
	IDTag            = "id"      // unique game id
	TypeTag          = "type"    // one of 'ai', 'multiplayer', or 'local'
	PlayersTagSuffix = "players" // tag suffix used when setting allowed players lists i.e. 'red_players': 'id1, id2'
)

const (
	AIType          = "ai"
	MultiplayerType = "multiplayer"
	LocalType       = "local"
	NoneType        = ""
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

func (t Tags) Players() (map[string][]string, error) {
	teams, err := t.Teams()
	if err != nil {
		return nil, err
	}
	typ, err := t.Type()
	if err != nil {
		return nil, err
	}

	players := make(map[string][]string)
	for _, team := range teams {
		tag := fmt.Sprintf("%s_%s", team, PlayersTagSuffix)
		p, ok := t[tag]
		if !ok {
			players[team] = []string{}
		} else {
			players[team] = strings.Split(p, ", ")
		}
	}

	if slices.Contains([]string{AIType}, typ) && len(players) != 1 {
		return nil, fmt.Errorf("players list must be of length one for ai games")
	}
	if typ == MultiplayerType && len(teams) != len(players) {
		return nil, fmt.Errorf("players list must match length of teams list for multiplayer games")
	}
	return players, nil
}

func (t Tags) Type() (string, error) {
	typ, ok := t[TypeTag]
	if !ok {
		return "", fmt.Errorf("optional %s tag is missing", TypeTag)
	}
	if !slices.Contains([]string{AIType, MultiplayerType, LocalType, NoneType}, typ) {
		return "", fmt.Errorf("invalid type tag value")
	}
	return typ, nil
}
