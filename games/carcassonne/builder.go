package carcassonne

import (
	"fmt"
	"strconv"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

const (
	Key = "carcassonne"
	Min = 2
	Max = 4
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
	game, err := NewCarcassonne(seed, teams)
	if err != nil {
		return nil, err
	}
	for _, action := range snapshot.Actions {
		if action.Index < 0 || action.Index >= len(teams) {
			return nil, fmt.Errorf("invalid action %d", action.Index)
		}
		team := teams[action.Index]
		switch action.Key {
		case ActionToQGN[PlaceTileAction]:
			if len(action.Details) != 9 {
				return nil, fmt.Errorf("invalid action details %v", action.Details)
			}
			x, err := strconv.Atoi(action.Details[0])
			if err != nil {
				return nil, err
			}
			y, err := strconv.Atoi(action.Details[1])
			if err != nil {
				return nil, err
			}
			top, ok := notationToStructure[action.Details[2]]
			if !ok {
				return nil, fmt.Errorf("failed to get top of tile")
			}
			right, ok := notationToStructure[action.Details[3]]
			if !ok {
				return nil, fmt.Errorf("failed to get right of tile")
			}
			bottom, ok := notationToStructure[action.Details[4]]
			if !ok {
				return nil, fmt.Errorf("failed to get bottom of tile")
			}
			left, ok := notationToStructure[action.Details[5]]
			if !ok {
				return nil, fmt.Errorf("failed to get left of tile")
			}
			center, ok := notationToStructure[action.Details[6]]
			if !ok {
				return nil, fmt.Errorf("failed to get center of tile")
			}
			connectedCitySides, ok := notationToBool[action.Details[7]]
			if !ok {
				return nil, fmt.Errorf("failed to get connectedCitySides")
			}
			banner, ok := notationToBool[action.Details[8]]
			if !ok {
				return nil, fmt.Errorf("failed to get banner")
			}
			if err := game.Do(&qg.Action{
				Team:    team,
				Type:    QGNToAction[action.Key],
				Details: PlaceTileDetails{x, y, Tile{top, right, bottom, left, center, connectedCitySides, banner}},
			}); err != nil {
				return nil, err
			}
		case ActionToQGN[PlaceTokenAction]:
			if len(action.Details) < 1 || len(action.Details) > 5 {
				return nil, fmt.Errorf("got %d but wanted %d to %d fields in when decoding %s details", len(action.Details), 1, 5, PlaceTokenAction)
			}
			pass, ok := notationToBool[action.Details[0]]
			if !ok {
				return nil, fmt.Errorf("got %s but wanted be 0 or 1 for for Pass when decoding %s details", action.Details[0], PlaceTokenAction)
			}
			if pass {
				if err := game.Do(&qg.Action{
					Team:    team,
					Type:    QGNToAction[action.Key],
					Details: PlaceTokenDetails{Pass: pass},
				}); err != nil {
					return nil, err
				}
			}
			x, err := strconv.Atoi(action.Details[1])
			if err != nil {
				return nil, err
			}
			y, err := strconv.Atoi(action.Details[2])
			if err != nil {
				return nil, err
			}
			token := notationToToken[action.Details[3]]
			if len(action.Details) == 4 {
				if err := game.Do(&qg.Action{
					Team:    team,
					Type:    QGNToAction[action.Key],
					Details: PlaceTokenDetails{Pass: pass, X: x, Y: y, Type: token},
				}); err != nil {
					return nil, err
				}
			}
			side := notationToSide[action.Details[4]]
			if token == Farmer {
				side = notationToFarmSide[action.Details[4]]
			}
			if err := game.Do(&qg.Action{
				Team:    team,
				Type:    QGNToAction[action.Key],
				Details: PlaceTokenDetails{Pass: pass, X: x, Y: y, Type: token, Side: side},
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
