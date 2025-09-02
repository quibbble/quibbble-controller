package tictactoe

import (
	"context"
	"fmt"
	"slices"
	"strings"

	sdk "github.com/quibbble/quibbble-controller/sdks/go"
	"google.golang.org/protobuf/types/known/emptypb"
	"k8s.io/utils/ptr"
)

const (
	Kind       = "tictactoe"
	MarkAction = "mark"
	NilTeam    = "nil"

	rows    = 3
	columns = 3
)

type TicTacToe struct {
	sdk.UnimplementedSDKServer
	*sdk.Snapshot
}

func (t *TicTacToe) InitGame(c context.Context, i *sdk.InitGameRequest) (*emptypb.Empty, error) {
	if i.Snapshot.Kind != Kind {
		return nil, fmt.Errorf("[%s init] %s is the wrong game", Kind, i.Snapshot.Kind)
	}
	t.Snapshot = i.Snapshot
	if len(t.History) == 0 {
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
		t.Message = ptr.To(fmt.Sprintf("%s must mark a spot", t.Turn))
	}
	return nil, nil
}

func (t *TicTacToe) DoAction(c context.Context, a *sdk.ActionRequest) (*emptypb.Empty, error) {
	switch a.Action.Kind {
	case MarkAction:
		markAction := TicTacToeMarkActionSpec{}
		if err := a.Action.Spec.UnmarshalTo(&markAction); err != nil {
			return nil, err
		}
		snapshotSpec := TicTacToeSnapshotSpec{}
		if err := t.Snapshot.Spec.UnmarshalTo(&snapshotSpec); err != nil {
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
		snapshotSpec.Row[markAction.Row].Column[markAction.Column] = a.Action.Team
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
			t.Snapshot.Winners = []string{winner}
			t.Message = ptr.To(fmt.Sprintf("%s wins", winner))
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

			t.Message = ptr.To(fmt.Sprintf("%s draw", strings.Join(t.Teams, ",")))
		}

		// change turn
		if len(t.Winners) == 0 {
			i := slices.Index(t.Teams, t.Turn)
			t.Turn = t.Teams[(i+1)%len(t.Teams)]
			t.Message = ptr.To(fmt.Sprintf("%s must mark a spot", t.Turn))
		}

		t.History = append(t.History, a.Action)
	}
	return nil, nil
}

func (t *TicTacToe) GetSnapshot(c context.Context, s *sdk.SnapshotRequest) (*sdk.SnapshotResponse, error) {
	return &sdk.SnapshotResponse{
		Snapshot: t.Snapshot,
	}, nil
}
