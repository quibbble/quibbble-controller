package tsuro

import (
	"fmt"
	"strconv"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

const (
	Key = "tsuro"
	Min = 2
	Max = 8
)

type Builder struct{}

func (b Builder) Create(snapshot *qgn.Snapshot) (qg.Game, error) {
	if key := snapshot.Tags[qgn.KeyTag]; key != Key {
		return nil, fmt.Errorf("%s is not a valid key", key)
	}
	teams, err := snapshot.Tags.Teams()
	if err != nil {
		return nil, err
	}
	variant, ok := snapshot.Tags[qgn.VariantTag]
	if !ok {
		variant = ClassicVariant
	}
	seed, err := snapshot.Tags.Seed()
	if err != nil {
		seed = time.Now().Unix()
	}
	game, err := NewTsuro(variant, seed, teams)
	if err != nil {
		return nil, err
	}
	for _, action := range snapshot.Actions {
		if action.Index < 0 || action.Index >= len(teams) {
			return nil, fmt.Errorf("invalid action %d", action.Index)
		}
		team := teams[action.Index]
		switch action.Key {
		case ActionToQGN[PlaceAction]:
			if len(action.Details) != 3 {
				return nil, fmt.Errorf("invalid action details %v", action.Details)
			}
			row, err := strconv.Atoi(action.Details[0])
			if err != nil {
				return nil, err
			}
			col, err := strconv.Atoi(action.Details[1])
			if err != nil {
				return nil, err
			}
			if err := game.Do(&qg.Action{
				Team:    team,
				Type:    QGNToAction[action.Key],
				Details: PlaceDetails{Row: row, Col: col, Tile: action.Details[2]},
			}); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid action key %s", action.Key)
		}
	}
	return game, nil
}

func (b Builder) GetInformation() *qg.Information {
	return &qg.Information{
		Key: Key,
		Min: Min,
		Max: Max,
	}
}
