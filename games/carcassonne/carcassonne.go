package carcassonne

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

type Carcassonne struct {
	*state
	history []*qg.Action
}

func NewCarcassonne(seed int64, teams []string) (*Carcassonne, error) {
	if len(teams) < Min || len(teams) > Max {
		return nil, fmt.Errorf("invalid number of teams")
	}

	return &Carcassonne{
		state:   newState(seed, teams),
		history: make([]*qg.Action, 0),
	}, nil
}

func (c *Carcassonne) Do(action *qg.Action) error {
	switch action.Type {
	case qg.ResetAction:
		g, err := qg.Reset(Builder{}, c, int(time.Now().Unix()))
		if err != nil {
			return err
		}
		c.state = g.(*Carcassonne).state
		c.history = make([]*qg.Action, 0)
	case qg.UndoAction:
		g, err := qg.Undo(Builder{}, c)
		if err != nil {
			return err
		}
		c.state = g.(*Carcassonne).state
		c.history = g.(*Carcassonne).history
	case qg.AIAction:
		if err := qg.AI(Builder{}, AI{}, c, 2); err != nil {
			return err
		}
	case RotateAction:
		if err := c.state.RotateTileRight(action.Team); err != nil {
			return err
		}
	case PlaceTileAction:
		var details PlaceTileDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		tile := newTile(details.Tile.Top, details.Tile.Right, details.Tile.Bottom, details.Tile.Left, details.Tile.Center, details.Tile.ConnectedCities, details.Tile.Banner)
		if err := c.PlaceTile(action.Team, tile, details.X, details.Y); err != nil {
			return err
		}
		c.history = append(c.history, action)
	case PlaceTokenAction:
		var details PlaceTokenDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := c.state.PlaceToken(action.Team, details.Pass, details.X, details.Y, details.Type, details.Side); err != nil {
			return err
		}
		c.history = append(c.history, action)
	default:
		return fmt.Errorf("%s is not a valid action", action.Type)
	}
	return nil
}

func (c *Carcassonne) GetSnapshotQGN() (*qgn.Snapshot, error) {
	tags := make(qgn.Tags)
	tags[qgn.KeyTag] = Key
	tags[qgn.TeamsTag] = strings.Join(c.teams, ", ")
	tags[qgn.SeedTag] = strconv.Itoa(int(c.seed))

	actions := make([]qgn.Action, 0)
	for _, action := range c.history {
		switch action.Type {
		case PlaceTileAction:
			var details PlaceTileDetails
			mapstructure.Decode(action.Details, &details)
			actions = append(actions, qgn.Action{
				Index: slices.Index(c.teams, action.Team),
				Key:   ActionToQGN[PlaceTileAction],
				Details: []string{strconv.Itoa(details.X), strconv.Itoa(details.Y),
					structureToNotation[details.Tile.Top], structureToNotation[details.Tile.Right],
					structureToNotation[details.Tile.Bottom], structureToNotation[details.Tile.Left],
					structureToNotation[details.Tile.Center],
					boolToNotation[details.Tile.ConnectedCities], boolToNotation[details.Tile.Banner]},
			})
		case PlaceTokenAction:
			var details PlaceTokenDetails
			mapstructure.Decode(action.Details, &details)
			var d []string
			if details.Pass {
				d = []string{boolToNotation[details.Pass]}
			} else if details.Type == Monk {
				d = []string{boolToNotation[details.Pass], strconv.Itoa(details.X), strconv.Itoa(details.Y), tokenToNotation[details.Type]}
			} else if details.Type == Farmer {
				d = []string{boolToNotation[details.Pass], strconv.Itoa(details.X), strconv.Itoa(details.Y), tokenToNotation[details.Type], farmSideToNotation[details.Side]}
			} else {
				d = []string{boolToNotation[details.Pass], strconv.Itoa(details.X), strconv.Itoa(details.Y), tokenToNotation[details.Type], sideToNotation[details.Side]}
			}
			actions = append(actions, qgn.Action{
				Index:   slices.Index(c.teams, action.Team),
				Key:     ActionToQGN[PlaceTokenAction],
				Details: d,
			})
		default:
			return nil, fmt.Errorf("%s is not a valid action", action.Type)
		}
	}

	return &qgn.Snapshot{
		Tags:    tags,
		Actions: actions,
	}, nil
}

func (c *Carcassonne) GetSnapshotJSON(team ...string) (*qg.Snapshot, error) {
	var actions []*qg.Action
	if len(c.winners) == 0 && (len(team) == 0 || (len(team) == 1 && team[0] == c.turn)) {
		actions = c.actions()
	}
	details := SnapshotDetails{
		LastPlaced:     c.state.lastPlacedTiles,
		Board:          c.state.board.board,
		BoardTokens:    c.state.boardTokens,
		Tokens:         c.state.tokens,
		Scores:         c.state.scores,
		TilesRemaining: len(c.state.deck.tiles),
	}
	if len(team) == 1 {
		details.PlayTile = c.state.playTiles[team[0]]
	}
	return &qg.Snapshot{
		Turn:    c.turn,
		Teams:   c.teams,
		Winners: c.winners,
		Details: details,
		Actions: actions,
		History: c.history,
		Message: c.message(),
	}, nil
}
