package onitama

import (
	"fmt"
	"strconv"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

const (
	Key = "onitama"
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
	seed, err := snapshot.Tags.Seed()
	if err != nil {
		seed = time.Now().Unix()
	}
	game, err := NewOnitama(seed, teams)
	if err != nil {
		return nil, err
	}
	for _, action := range snapshot.Actions {
		if action.Index < 0 || action.Index >= len(teams) {
			return nil, fmt.Errorf("invalid action %d", action.Index)
		}
		team := teams[action.Index]
		switch action.Key {
		case ActionToQGN[MoveAction]:
			if len(action.Details) != 5 {
				return nil, fmt.Errorf("invalid action details %v", action.Details)
			}
			startRow, err := strconv.Atoi(action.Details[0])
			if err != nil {
				return nil, err
			}
			startCol, err := strconv.Atoi(action.Details[1])
			if err != nil {
				return nil, err
			}
			endRow, err := strconv.Atoi(action.Details[2])
			if err != nil {
				return nil, err
			}
			endCol, err := strconv.Atoi(action.Details[3])
			if err != nil {
				return nil, err
			}
			card := action.Details[4]
			if err := game.Do(&qg.Action{
				Team:    team,
				Type:    QGNToAction[action.Key],
				Details: MoveDetails{StartRow: startRow, StartCol: startCol, EndRow: endRow, EndCol: endCol, Card: card},
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
