package tictactoe

import (
	"encoding/base64"
	"testing"

	q "github.com/quibbble/quibbble-controller/pkg/quibbble"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func Test_TicTactToe(t *testing.T) {
	snapSpec, err := anypb.New(&TicTacToeSnapshotSpec{
		Row: []*TicTacToeRow{
			{Column: []string{NilTeam, NilTeam, NilTeam}},
			{Column: []string{NilTeam, NilTeam, NilTeam}},
			{Column: []string{NilTeam, NilTeam, NilTeam}},
		},
	})
	if err != nil {
		t.FailNow()
	}
	snap := q.GameSnapshot{
		Turn:  "X",
		Teams: []string{"X", "O"},
		Spec:  snapSpec,
	}

	println(snap.String())

	raw, err := proto.Marshal(&snap)
	if err != nil {
		t.FailNow()
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(raw))

	println(b64)

	raw2, err := base64.RawStdEncoding.DecodeString(b64)
	if err != nil {
		t.FailNow()
	}
	snap2 := q.GameSnapshot{}
	if err := proto.Unmarshal(raw2, &snap2); err != nil {
		t.FailNow()
	}

	println(snap2.String())
}
