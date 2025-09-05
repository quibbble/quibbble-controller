package tictactoe

import (
	"context"
	"fmt"
	"slices"

	q "github.com/quibbble/quibbble-controller/pkg/quibbble"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	Kind       = "tictactoe"
	MarkAction = "mark"
	NilTeam    = "nil"

	rows    = 3
	columns = 3
)

type TicTacToe struct {
	q.UnimplementedGameServer
	*q.GameSnapshot
}

func validateSnapshot(s *q.GameSnapshot) error {
	if len(s.Teams) != 2 {
		return fmt.Errorf("only two teams are allow")
	}
	if s.Teams[0] == s.Teams[1] {
		return fmt.Errorf("teams cannot be the same")
	}
	if !slices.Contains(s.Teams, s.Turn) {
		return fmt.Errorf("%s is not a valid team", s.Turn)
	}
	return nil
}

func (t *TicTacToe) Init(c context.Context, s *q.GameSnapshot) (*emptypb.Empty, error) {
	if err := validateSnapshot(s); err != nil {
		return nil, err
	}
	t.GameSnapshot = &q.GameSnapshot{
		Teams: s.Teams,
		Turn:  s.Turn,
	}
	snapshotSpec := TicTacToeSnapshotSpec{
		Row: []*TicTacToeRow{
			{Column: []string{NilTeam, NilTeam, NilTeam}},
			{Column: []string{NilTeam, NilTeam, NilTeam}},
			{Column: []string{NilTeam, NilTeam, NilTeam}},
		},
	}
	if err := t.Spec.MarshalFrom(&snapshotSpec); err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *TicTacToe) Load(c context.Context, s *q.GameSnapshot) (*emptypb.Empty, error) {
	if err := validateSnapshot(s); err != nil {
		return nil, err
	}
	t.GameSnapshot = s
	return nil, nil
}

func (t *TicTacToe) PlayAction(c context.Context, a *q.GameAction) (*emptypb.Empty, error) {
	switch a.Kind {
	case MarkAction:
		markAction := TicTacToeMarkActionSpec{}
		if err := a.Spec.UnmarshalTo(&markAction); err != nil {
			return nil, err
		}
		snapshotSpec := TicTacToeSnapshotSpec{}
		if err := t.GameSnapshot.Spec.UnmarshalTo(&snapshotSpec); err != nil {
			return nil, err
		}

		// check if mark action can be taken
		if len(t.Winners) > 0 {
			return nil, fmt.Errorf("[%s %s] game is already over", Kind, MarkAction)
		}
		if markAction.Row < 0 || markAction.Column < 0 ||
			markAction.Row > rows || markAction.Column > columns {
			return nil, fmt.Errorf("[%s %s] row or column out of bounds", Kind, MarkAction)
		}
		if snapshotSpec.Row[markAction.Row].Column[markAction.Column] != NilTeam {
			return nil, fmt.Errorf("[%s %s] location is already taken", Kind, MarkAction)
		}

		// do mark action
		snapshotSpec.Row[markAction.Row].Column[markAction.Column] = a.Team
		if err := t.Spec.MarshalFrom(&snapshotSpec); err != nil {
			return nil, err
		}

		// check for winner
		equal := func(a, b, c string) bool {
			if a != NilTeam && a == b && b == c {
				return true
			}
			return false
		}

		winner := NilTeam
		for i := range rows {
			if equal(snapshotSpec.Row[0].Column[i], snapshotSpec.Row[1].Column[i], snapshotSpec.Row[2].Column[i]) {
				winner = snapshotSpec.Row[0].Column[i]
			}
			if equal(snapshotSpec.Row[i].Column[0], snapshotSpec.Row[i].Column[1], snapshotSpec.Row[i].Column[2]) {
				winner = snapshotSpec.Row[i].Column[0]
			}
		}
		if equal(snapshotSpec.Row[0].Column[0], snapshotSpec.Row[1].Column[1], snapshotSpec.Row[2].Column[2]) {
			winner = snapshotSpec.Row[0].Column[0]
		}
		if equal(snapshotSpec.Row[2].Column[0], snapshotSpec.Row[1].Column[1], snapshotSpec.Row[0].Column[2]) {
			winner = snapshotSpec.Row[2].Column[0]
		}
		if winner != NilTeam {
			t.GameSnapshot.Winners = []string{winner}
		}

		// check for draw
		draw := true
		for i := range rows {
			for j := range columns {
				if snapshotSpec.Row[i].Column[j] == NilTeam {
					draw = false
				}
			}
		}
		if draw {
			t.Winners = t.Teams
		}

		// change turn
		if len(t.Winners) == 0 {
			i := slices.Index(t.Teams, t.Turn)
			t.Turn = t.Teams[(i+1)%len(t.Teams)]
		}
	}
	return nil, nil
}

func (t *TicTacToe) GetSnapshot(c context.Context, s *q.GameView) (*q.GameSnapshot, error) {
	return t.GameSnapshot, nil
}
