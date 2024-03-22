package connect4

import "github.com/quibbble/quibbble-controller/pkg/util"

const (
	Rows = 6
	Cols = 7
)

const PlaceAction = "place"

var Actions = []string{PlaceAction}

var (
	ActionToQGN = map[string]string{
		PlaceAction: "p",
	}
	QGNToAction = util.ReverseMap(ActionToQGN)
)

type PlaceDetails struct {
	Col int `json:"col"`
}

type SnapshotDetails struct {
	Board [Rows][Cols]*string `json:"board"`
}
