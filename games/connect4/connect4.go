package connect4

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

type Connect4 struct {
	*state
	history []*qg.Action
}

func NewConnect4(teams []string) (*Connect4, error) {
	if len(teams) < Min || len(teams) > Max {
		return nil, fmt.Errorf("invalid number of teams")
	}
	return &Connect4{
		state:   newState(teams),
		history: make([]*qg.Action, 0),
	}, nil
}

func (c *Connect4) Do(action *qg.Action) error {
	switch action.Type {
	case qg.ResetAction:
		g, err := qg.Reset(Builder{}, c)
		if err != nil {
			return err
		}
		c.state = g.(*Connect4).state
		c.history = make([]*qg.Action, 0)
	case qg.UndoAction:
		g, err := qg.Undo(Builder{}, c)
		if err != nil {
			return err
		}
		c.state = g.(*Connect4).state
		c.history = g.(*Connect4).history
	case qg.AIAction:
		if err := qg.AI(Builder{}, AI{}, c, 3); err != nil {
			return err
		}
	case PlaceAction:
		var details PlaceDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := c.Place(action.Team, details.Col); err != nil {
			return err
		}
		c.history = append(c.history, action)
	default:
		return fmt.Errorf("%s is not a valid action", action.Type)
	}
	return nil
}

func (c *Connect4) GetSnapshotQGN() (*qgn.Snapshot, error) {
	tags := make(qgn.Tags)
	tags[qgn.KeyTag] = Key
	tags[qgn.TeamsTag] = strings.Join(c.teams, ", ")

	actions := make([]qgn.Action, 0)
	for _, action := range c.history {
		switch action.Type {
		case PlaceAction:
			var details PlaceDetails
			mapstructure.Decode(action.Details, &details)
			actions = append(actions, qgn.Action{
				Index:   slices.Index(c.teams, action.Team),
				Key:     ActionToQGN[PlaceAction],
				Details: []string{strconv.Itoa(details.Col)},
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

func (c *Connect4) GetSnapshotJSON(team ...string) (*qg.Snapshot, error) {
	var actions []*qg.Action
	if len(c.winners) == 0 && (len(team) == 0 || (len(team) == 1 && team[0] == c.turn)) {
		actions = c.actions()
	}
	return &qg.Snapshot{
		Turn:    c.turn,
		Teams:   c.teams,
		Winners: c.winners,
		Details: SnapshotDetails{
			Board: c.board.board,
		},
		Actions: actions,
		History: c.history,
		Message: c.message(),
	}, nil
}
