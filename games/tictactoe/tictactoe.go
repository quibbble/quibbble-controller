package tictactoe

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

type Tictactoe struct {
	*state
	history []*qg.Action
}

func NewTicTacToe(teams []string) (*Tictactoe, error) {
	if len(teams) < Min || len(teams) > Max {
		return nil, fmt.Errorf("invalid number of teams")
	}
	return &Tictactoe{
		state: &state{
			turn:    teams[0],
			teams:   teams,
			winners: make([]string, 0),
			board:   [BoardSize][BoardSize]*string{},
		},
		history: make([]*qg.Action, 0),
	}, nil
}

func (t *Tictactoe) Do(action *qg.Action) error {
	switch action.Type {
	case qg.ResetAction:
		g, err := qg.Reset(Builder{}, t)
		if err != nil {
			return err
		}
		t.state = g.(*Tictactoe).state
		t.history = make([]*qg.Action, 0)
	case qg.UndoAction:
		g, err := qg.Undo(Builder{}, t)
		if err != nil {
			return err
		}
		t.state = g.(*Tictactoe).state
		t.history = g.(*Tictactoe).history
	case qg.AIAction:
		if err := qg.AI(Builder{}, AI{}, t, 3); err != nil {
			return err
		}
	case MarkAction:
		var details MarkDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := t.Mark(action.Team, details.Row, details.Col); err != nil {
			return err
		}
		t.history = append(t.history, action)
	default:
		return fmt.Errorf("%s is not a valid action", action.Type)
	}
	return nil
}

func (t *Tictactoe) GetSnapshotQGN() (*qgn.Snapshot, error) {
	tags := make(qgn.Tags)
	tags[qgn.KeyTag] = Key
	tags[qgn.TeamsTag] = strings.Join(t.teams, ", ")

	actions := make([]qgn.Action, 0)
	for _, action := range t.history {
		switch action.Type {
		case MarkAction:
			var details MarkDetails
			mapstructure.Decode(action.Details, &details)
			actions = append(actions, qgn.Action{
				Index:   slices.Index(t.teams, action.Team),
				Key:     ActionToQGN[MarkAction],
				Details: []string{strconv.Itoa(details.Row), strconv.Itoa(details.Col)},
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

func (t *Tictactoe) GetSnapshotJSON(team ...string) (*qg.Snapshot, error) {
	var actions []*qg.Action
	if len(t.winners) == 0 && (len(team) == 0 || (len(team) == 1 && team[0] == t.turn)) {
		actions = t.actions()
	}
	return &qg.Snapshot{
		Turn:    t.turn,
		Teams:   t.teams,
		Winners: t.winners,
		Details: SnapshotDetails{
			Board: t.board,
		},
		Actions: actions,
		History: t.history,
		Message: t.message(),
	}, nil
}
