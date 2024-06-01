package stratego

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

type Stratego struct {
	*state
	history []*qg.Action
}

func NewStratego(variant string, seed int64, teams []string) (*Stratego, error) {
	if len(teams) < Min || len(teams) > Max {
		return nil, fmt.Errorf("invalid number of teams")
	}
	state, err := newState(variant, seed, teams)
	if err != nil {
		return nil, err
	}
	return &Stratego{
		state:   state,
		history: make([]*qg.Action, 0),
	}, nil
}

func (s *Stratego) Do(action *qg.Action) error {
	switch action.Type {
	case qg.ResetAction:
		g, err := qg.Reset(Builder{}, s, int(time.Now().Unix()))
		if err != nil {
			return err
		}
		s.state = g.(*Stratego).state
		s.history = make([]*qg.Action, 0)
	case qg.UndoAction:
		g, err := qg.Undo(Builder{}, s)
		if err != nil {
			return err
		}
		s.state = g.(*Stratego).state
		s.history = g.(*Stratego).history
	case qg.AIAction:
		if err := qg.AI(Builder{}, AI{}, s, 2); err != nil {
			return err
		}
	case SwitchAction:
		var details SwitchDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := s.Switch(action.Team, details.UnitARow, details.UnitACol, details.UnitBRow, details.UnitBCol); err != nil {
			return err
		}
		s.history = append(s.history, action)
	case MoveAction:
		var details MoveDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := s.Move(action.Team, details.UnitRow, details.UnitCol, details.TileRow, details.TileCol); err != nil {
			return err
		}
		s.history = append(s.history, action)
	default:
		return fmt.Errorf("%s is not a valid action", action.Type)
	}
	return nil
}

func (s *Stratego) GetSnapshotQGN() (*qgn.Snapshot, error) {
	tags := make(qgn.Tags)
	tags[qgn.KeyTag] = Key
	tags[qgn.TeamsTag] = strings.Join(s.teams, ", ")
	tags[qgn.SeedTag] = strconv.Itoa(int(s.seed))

	actions := make([]qgn.Action, 0)
	for _, action := range s.history {
		switch action.Type {
		case SwitchAction:
			var details SwitchDetails
			mapstructure.Decode(action.Details, &details)
			actions = append(actions, qgn.Action{
				Index: slices.Index(s.teams, action.Team),
				Key:   ActionToQGN[SwitchAction],
				Details: []string{strconv.Itoa(details.UnitARow), strconv.Itoa(details.UnitACol),
					strconv.Itoa(details.UnitBRow), strconv.Itoa(details.UnitBCol)},
			})
		case MoveAction:
			var details MoveDetails
			mapstructure.Decode(action.Details, &details)
			actions = append(actions, qgn.Action{
				Index: slices.Index(s.teams, action.Team),
				Key:   ActionToQGN[MoveAction],
				Details: []string{strconv.Itoa(details.UnitRow), strconv.Itoa(details.UnitCol),
					strconv.Itoa(details.TileRow), strconv.Itoa(details.TileCol)},
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

func (s *Stratego) GetSnapshotJSON(team ...string) (*qg.Snapshot, error) {
	var actions []*qg.Action
	if len(s.winners) == 0 && (len(team) == 0 || (len(team) == 1 && team[0] == s.turn)) {
		actions = s.actions()
	}

	// reveals the winning unit from the last battle to both teams
	revealRow := -1
	revealCol := -1
	if s.state.battle != nil && s.state.justBattled {
		revealRow = s.state.battle.TileRow
		revealCol = s.state.battle.TileCol
	}

	board := [][]Unit{}
	for r, row := range s.state.board.board {
		sRow := make([]Unit, 0)
		for c, unit := range row {
			if unit == nil {
				sRow = append(sRow, Unit{})
			} else {
				if unit.Team != nil {
					if len(team) == 1 {
						if *unit.Team == team[0] || revealRow == r && revealCol == c || len(s.state.winners) > 0 {
							sRow = append(sRow, *NewUnit(unit.Type, *unit.Team))
						} else {
							sRow = append(sRow, *NewUnit("", *unit.Team))
						}
					} else {
						sRow = append(sRow, *NewUnit("", *unit.Team))
					}
				} else {
					sRow = append(sRow, *Water())
				}
			}
		}
		board = append(board, sRow)
	}

	return &qg.Snapshot{
		Turn:    s.state.turn,
		Teams:   s.teams,
		Winners: s.winners,
		Details: SnapshotDetails{
			Board:       board,
			Battle:      s.state.battle,
			JustBattled: s.state.justBattled,
			Started:     s.state.started,
			Variant:     s.variant,
		},
		Actions: actions,
		History: s.history,
		Message: s.message(),
	}, nil
}
