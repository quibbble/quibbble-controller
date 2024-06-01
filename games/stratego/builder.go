package stratego

import (
	"fmt"
	"strconv"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

const (
	Key = "stratego"
	Min = 2
	Max = 2
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
	game, err := NewStratego(variant, seed, teams)
	if err != nil {
		return nil, err
	}
	for _, action := range snapshot.Actions {
		if action.Index < 0 || action.Index >= len(teams) {
			return nil, fmt.Errorf("invalid action %d", action.Index)
		}
		team := teams[action.Index]
		switch action.Key {
		case ActionToQGN[SwitchAction]:
			if len(action.Details) != 4 {
				return nil, fmt.Errorf("invalid action details %v", action.Details)
			}
			unitARow, err := strconv.Atoi(action.Details[0])
			if err != nil {
				return nil, err
			}
			unitACol, err := strconv.Atoi(action.Details[1])
			if err != nil {
				return nil, err
			}
			unitBRow, err := strconv.Atoi(action.Details[2])
			if err != nil {
				return nil, err
			}
			unitBCol, err := strconv.Atoi(action.Details[3])
			if err != nil {
				return nil, err
			}
			if err := game.Do(&qg.Action{
				Team:    team,
				Type:    QGNToAction[action.Key],
				Details: SwitchDetails{unitARow, unitACol, unitBRow, unitBCol},
			}); err != nil {
				return nil, err
			}
		case ActionToQGN[MoveAction]:
			if len(action.Details) != 4 {
				return nil, fmt.Errorf("invalid action details %v", action.Details)
			}
			unitRow, err := strconv.Atoi(action.Details[0])
			if err != nil {
				return nil, err
			}
			unitCol, err := strconv.Atoi(action.Details[1])
			if err != nil {
				return nil, err
			}
			tileRow, err := strconv.Atoi(action.Details[2])
			if err != nil {
				return nil, err
			}
			tileCol, err := strconv.Atoi(action.Details[3])
			if err != nil {
				return nil, err
			}
			if err := game.Do(&qg.Action{
				Team:    team,
				Type:    QGNToAction[action.Key],
				Details: MoveDetails{unitRow, unitCol, tileRow, tileCol},
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
