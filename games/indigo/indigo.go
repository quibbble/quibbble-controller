package indigo

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

type Indigo struct {
	*state
	history []*qg.Action
}

func NewIndigo(variant string, seed int64, teams []string) (*Indigo, error) {
	if len(teams) < Min || len(teams) > Max {
		return nil, fmt.Errorf("invalid number of teams")
	}
	state, err := newState(variant, seed, teams)
	if err != nil {
		return nil, err
	}
	return &Indigo{
		state:   state,
		history: make([]*qg.Action, 0),
	}, nil
}

func (i *Indigo) Do(action *qg.Action) error {
	switch action.Type {
	case qg.ResetAction:
		g, err := qg.Reset(Builder{}, i, int(time.Now().Unix()))
		if err != nil {
			return err
		}
		i.state = g.(*Indigo).state
		i.history = make([]*qg.Action, 0)
	case qg.UndoAction:
		g, err := qg.Undo(Builder{}, i)
		if err != nil {
			return err
		}
		i.state = g.(*Indigo).state
		i.history = g.(*Indigo).history
	case qg.AIAction:
		if err := qg.AI(Builder{}, AI{}, i, 1); err != nil {
			return err
		}
	case PlaceAction:
		var details PlaceDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := i.place(action.Team, details.Tile, details.Row, details.Col); err != nil {
			return err
		}
		i.history = append(i.history, action)
	case RotateAction:
		var details RotateDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := i.rotate(action.Team, details.Tile); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s is not a valid action", action.Type)
	}
	return nil
}

func (i *Indigo) GetSnapshotQGN() (*qgn.Snapshot, error) {
	tags := make(qgn.Tags)
	tags[qgn.KeyTag] = Key
	tags[qgn.TeamsTag] = strings.Join(i.teams, ", ")
	tags[qgn.SeedTag] = strconv.Itoa(int(i.seed))

	actions := make([]qgn.Action, 0)
	for _, action := range i.history {
		switch action.Type {
		case PlaceAction:
			var details PlaceDetails
			mapstructure.Decode(action.Details, &details)
			actions = append(actions, qgn.Action{
				Index:   slices.Index(i.teams, action.Team),
				Key:     ActionToQGN[PlaceAction],
				Details: []string{details.Tile, strconv.Itoa(details.Row), strconv.Itoa(details.Col)},
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

func (i *Indigo) GetSnapshotJSON(team ...string) (*qg.Snapshot, error) {
	var actions []*qg.Action
	if len(i.winners) == 0 && (len(team) == 0 || (len(team) == 1 && team[0] == i.turn)) {
		actions = i.actions(team...)
	}

	hands := make(map[string][]tile)
	for t, hand := range i.state.hands {
		if len(team) == 0 || (t == team[0]) {
			hands[t] = hand.GetItems()
		}
	}

	return &qg.Snapshot{
		Turn:    i.turn,
		Teams:   i.teams,
		Winners: i.winners,
		Details: SnapshotDetails{
			Board:          i.state.board,
			Hands:          hands,
			Points:         i.state.points,
			Round:          i.state.round,
			RoundsUntilEnd: i.state.roundsUntilEnd,
			Variant:        i.state.variant,
		},
		Actions: actions,
		History: i.history,
		Message: i.message(),
	}, nil
}
