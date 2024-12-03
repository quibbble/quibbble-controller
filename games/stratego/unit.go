package stratego

import (
	"fmt"
	"slices"
)

// Unit Types
const (
	water = "water"

	flag       = "flag"
	bomb       = "bomb"
	spy        = "spy"
	scout      = "scout"
	miner      = "miner"
	sergeant   = "sergeant"
	lieutenant = "lieutenant"
	captain    = "captain"
	major      = "major"
	colonel    = "colonel"
	general    = "general"
	marshal    = "marshal"
)

// UnitTyes ordered in ascending order by battle winner
var UnitTyes = []string{flag, bomb, spy, scout, miner, sergeant, lieutenant, captain, major, colonel, general, marshal}

type Unit struct {
	Type string  `json:"type"`
	Team *string `json:"team,omitempty"`
}

func NewUnit(typ, team string) *Unit {
	return &Unit{
		Type: typ,
		Team: &team,
	}
}

func Water() *Unit {
	return &Unit{
		Type: water,
		Team: nil,
	}
}

func (u *Unit) Attack(unit *Unit) (winner *Unit, err error) {
	if *u.Team == *unit.Team {
		return nil, fmt.Errorf("cannot attack unit on same team")
	}
	if u.Type == flag || u.Type == bomb {
		return nil, fmt.Errorf("%s cannot attack", u.Type)
	}
	// spy -> marshal case
	if u.Type == spy && unit.Type == marshal {
		return u, nil
	}
	// miner -> bomb case
	if u.Type == miner && unit.Type == bomb {
		return u, nil
	}
	// any -> bomb case
	if unit.Type == bomb {
		return unit, nil
	}
	// same type case
	if u.Type == unit.Type {
		return nil, nil
	}
	// default case
	if slices.Index(UnitTyes, u.Type) > slices.Index(UnitTyes, unit.Type) {
		return u, nil
	} else {
		return unit, nil
	}
}
