package tictactoe

import "github.com/quibbble/quibbble-controller/pkg/util"

const BoardSize = 3

const MarkAction = "mark"

var AllActions = []string{MarkAction}

var (
	ActionToQGN = map[string]string{
		MarkAction: "m",
	}
	QGNToAction = util.ReverseMap(ActionToQGN)
)

type MarkDetails struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type SnapshotDetails struct {
	Board [BoardSize][BoardSize]*string `json:"board"`
}
