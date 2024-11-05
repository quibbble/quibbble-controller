package onitama

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

type Onitama struct {
	*state
	history []*qg.Action
}

func NewOnitama(seed int64, teams []string) (*Onitama, error) {
	if len(teams) < Min || len(teams) > Max {
		return nil, fmt.Errorf("invalid number of teams")
	}

	// select the 5 cards used for this game
	r := rand.New(rand.NewSource(seed))
	selected := make([]string, 0)
	for len(selected) < 5 {
		i := r.Intn(len(cards))
		if !slices.Contains(selected, cards[i]) {
			selected = append(selected, cards[i])
		}
	}

	hands := make(map[string][]string)
	hands[teams[0]] = []string{selected[0], selected[1]}
	hands[teams[1]] = []string{selected[2], selected[3]}

	board := [BoardSize][BoardSize]*Pawn{}
	board[0][0] = &Pawn{Type: student, Team: teams[0]}
	board[0][1] = &Pawn{Type: student, Team: teams[0]}
	board[0][2] = &Pawn{Type: master, Team: teams[0]}
	board[0][3] = &Pawn{Type: student, Team: teams[0]}
	board[0][4] = &Pawn{Type: student, Team: teams[0]}
	board[BoardSize-1][0] = &Pawn{Type: student, Team: teams[1]}
	board[BoardSize-1][1] = &Pawn{Type: student, Team: teams[1]}
	board[BoardSize-1][2] = &Pawn{Type: master, Team: teams[1]}
	board[BoardSize-1][3] = &Pawn{Type: student, Team: teams[1]}
	board[BoardSize-1][4] = &Pawn{Type: student, Team: teams[1]}

	return &Onitama{
		state: &state{
			seed:    seed,
			turn:    teams[0],
			teams:   teams,
			winners: make([]string, 0),
			board:   board,
			hands:   hands,
			spare:   selected[4],
		},
		history: make([]*qg.Action, 0),
	}, nil
}

func (o *Onitama) Do(action *qg.Action) error {
	switch action.Type {
	case qg.ResetAction:
		g, err := qg.Reset(Builder{}, o)
		if err != nil {
			return err
		}
		o.state = g.(*Onitama).state
		o.history = make([]*qg.Action, 0)
	case qg.UndoAction:
		g, err := qg.Undo(Builder{}, o)
		if err != nil {
			return err
		}
		o.state = g.(*Onitama).state
		o.history = g.(*Onitama).history
	case qg.AIAction:
		if err := qg.AI(Builder{}, AI{}, o, 3); err != nil {
			return err
		}
	case MoveAction:
		var details MoveDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := o.Move(action.Team, details.StartRow, details.StartCol, details.EndRow, details.EndCol, details.Card); err != nil {
			return err
		}
		o.history = append(o.history, action)
	default:
		return fmt.Errorf("%s is not a valid action", action.Type)
	}
	return nil
}

func (o *Onitama) GetSnapshotQGN() (*qgn.Snapshot, error) {
	tags := make(qgn.Tags)
	tags[qgn.KeyTag] = Key
	tags[qgn.TeamsTag] = strings.Join(o.teams, ", ")
	tags[qgn.SeedTag] = strconv.Itoa(int(o.seed))

	actions := make([]qgn.Action, 0)
	for _, action := range o.history {
		switch action.Type {
		case MoveAction:
			var details MoveDetails
			mapstructure.Decode(action.Details, &details)
			actions = append(actions, qgn.Action{
				Index:   slices.Index(o.teams, action.Team),
				Key:     ActionToQGN[MoveAction],
				Details: []string{strconv.Itoa(details.StartRow), strconv.Itoa(details.StartCol), strconv.Itoa(details.EndRow), strconv.Itoa(details.EndCol), details.Card},
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

func (o *Onitama) GetSnapshotJSON(team ...string) (*qg.Snapshot, error) {
	var actions []*qg.Action
	if len(o.winners) == 0 && (len(team) == 0 || (len(team) == 1 && team[0] == o.turn)) {
		actions = o.actions()
	}
	return &qg.Snapshot{
		Turn:    o.turn,
		Teams:   o.teams,
		Winners: o.winners,
		Details: SnapshotDetails{
			Board: o.board,
			Hands: o.hands,
			Spare: o.spare,
		},
		Actions: actions,
		History: o.history,
		Message: o.message(),
	}, nil
}
