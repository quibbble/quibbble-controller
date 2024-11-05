package onitama

import "github.com/quibbble/quibbble-controller/pkg/util"

const BoardSize = 5

const MoveAction = "move"

var Actions = []string{MoveAction}

var (
	ActionToQGN = map[string]string{
		MoveAction: "m",
	}
	QGNToAction = util.ReverseMap(ActionToQGN)
)

type MoveDetails struct {
	StartRow int `json:"start_row"`
	StartCol int `json:"start_col"`

	EndRow int `json:"end_row"`
	EndCol int `json:"end_col"`

	Card string `json:"card"`
}

type SnapshotDetails struct {
	Board [BoardSize][BoardSize]*Pawn `json:"board"`
	Hands map[string][]string         `json:"hands"`
	Spare string                      `json:"spare"`
}

type Pawn struct {
	Type string `json:"type"`
	Team string `json:"team"`
}

// pawn types
const (
	master  = "master"
	student = "student"
)

// card types
const (
	tiger    = "tiger"
	dragon   = "dragon"
	frog     = "frog"
	rabbit   = "rabbit"
	crab     = "crab"
	elephant = "elephant"
	goose    = "goose"
	rooster  = "rooster"
	monkey   = "monkey"
	mantis   = "mantis"
	horse    = "horse"
	ox       = "ox"
	crane    = "crane"
	boar     = "boar"
	eel      = "eel"
	cobra    = "cobra"
)

var cards = []string{tiger, dragon, frog, rabbit, crab, elephant, goose, rooster, monkey, mantis, horse, ox, crane, boar, eel, cobra}

// map from card to valid (x, y) moves using the card
var movements = map[string][][2]int{
	tiger:    {{0, 2}, {0, -1}},
	dragon:   {{2, 1}, {-2, 1}, {1, -1}, {-1, -1}},
	frog:     {{1, -1}, {-1, 1}, {-2, 0}},
	rabbit:   {{-1, -1}, {1, 1}, {2, 0}},
	crab:     {{0, 1}, {2, 0}, {-2, 0}},
	elephant: {{1, 0}, {-1, 0}, {1, 1}, {-1, 1}},
	goose:    {{1, 0}, {-1, 0}, {1, -1}, {-1, 1}},
	rooster:  {{1, 0}, {-1, 0}, {1, 1}, {-1, -1}},
	monkey:   {{1, 1}, {-1, -1}, {-1, 1}, {1, -1}},
	mantis:   {{1, 1}, {-1, 1}, {0, -1}},
	horse:    {{0, 1}, {0, -1}, {-1, 0}},
	ox:       {{0, 1}, {0, -1}, {1, 0}},
	crane:    {{0, 1}, {1, -1}, {-1, -1}},
	boar:     {{0, 1}, {1, 0}, {-1, 0}},
	eel:      {{-1, 1}, {1, 0}, {-1, -1}},
	cobra:    {{-1, 0}, {1, 1}, {1, -1}},
}
