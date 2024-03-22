package tsuro

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

type Tsuro struct {
	*state
	history []*qg.Action
}

func NewTsuro(variant string, seed int64, teams []string) (*Tsuro, error) {
	if len(teams) < Min || len(teams) > Max {
		return nil, fmt.Errorf("invalid number of teams")
	}
	if !slices.Contains(Variants, variant) {
		return nil, fmt.Errorf("variant %s is not supported", variant)
	}
	state, err := newState(variant, seed, teams)
	if err != nil {
		return nil, err
	}
	return &Tsuro{
		state:   state,
		history: make([]*qg.Action, 0),
	}, nil
}

func (t *Tsuro) Do(action *qg.Action) error {
	switch action.Type {
	case qg.ResetAction:
		g, err := qg.Reset(Builder{}, t)
		if err != nil {
			return err
		}
		t.state = g.(*Tsuro).state
		t.history = make([]*qg.Action, 0)
	case qg.UndoAction:
		g, err := qg.Undo(Builder{}, t)
		if err != nil {
			return err
		}
		t.state = g.(*Tsuro).state
		t.history = g.(*Tsuro).history
	case qg.AIAction:
		if err := qg.AI(Builder{}, AI{}, t, 3); err != nil {
			return err
		}
	case PlaceAction:
		var details PlaceDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := t.Place(action.Team, details.Tile, details.Row, details.Col); err != nil {
			return err
		}
		t.history = append(t.history, action)
	case RotateAction:
		var details RotateDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := t.Rotate(action.Team, details.Tile); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s is not a valid action", action.Type)
	}
	return nil
}

func (t *Tsuro) GetSnapshotQGN() (*qgn.Snapshot, error) {
	tags := make(qgn.Tags)
	tags[qgn.KeyTag] = Key
	tags[qgn.VariantTag] = t.variant
	tags[qgn.SeedTag] = strconv.Itoa(int(t.seed))
	tags[qgn.TeamsTag] = strings.Join(t.teams, ", ")

	actions := make([]qgn.Action, 0)
	for _, action := range t.history {
		switch action.Type {
		case PlaceAction:
			var details PlaceDetails
			mapstructure.Decode(action.Details, &details)
			actions = append(actions, qgn.Action{
				Index:   slices.Index(t.teams, action.Team),
				Key:     ActionToQGN[PlaceAction],
				Details: []string{strconv.Itoa(details.Row), strconv.Itoa(details.Col), details.Tile},
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

func (t *Tsuro) GetSnapshotJSON(team ...string) (*qg.Snapshot, error) {
	hands := make(map[string][]*tile)
	for t, hand := range t.state.hands {
		if len(team) == 0 {
			hands[t] = hand.hand
		} else {
			if team[0] == t {
				hands[t] = hand.hand
			}
		}
	}
	var points map[string]int
	if t.state.variant == LongestPathVariant || t.state.variant == MostCrossingsVariant {
		points = t.state.points
	}
	details := SnapshotDetails{
		Board:          t.state.board.board,
		TilesRemaining: len(t.state.deck.deck),
		Hands:          hands,
		Tokens:         t.state.tokens,
		Dragon:         t.state.dragon,
		Variant:        t.state.variant,
		Points:         points,
	}
	var actions []*qg.Action
	if len(t.state.winners) == 0 {
		actions = t.state.actions(team...)
	}
	return &qg.Snapshot{
		Turn:    t.state.turn,
		Teams:   t.state.teams,
		Winners: t.state.winners,
		Details: details,
		Actions: actions,
		History: t.history,
		Message: t.state.message(),
	}, nil
}
