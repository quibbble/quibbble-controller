package indigo

import "github.com/quibbble/quibbble-controller/pkg/util"

// Action types
const (
	RotateAction = "rotate" // NOTE - this is not tracked by BGN
	PlaceAction  = "place"
)

var Actions = []string{RotateAction, PlaceAction}

var (
	ActionToQGN = map[string]string{
		PlaceAction: "p",
	}
	QGNToAction = util.ReverseMap(ActionToQGN)
)

const (
	ClassicVariant    = "classic"     // normal Indigo
	LargeHandsVariant = "large_hands" // players have a hand size of 2 instead of 1
)

var Variants = []string{ClassicVariant, LargeHandsVariant}

type RotateDetails struct {
	Tile string `json:"tile"`
}

type PlaceDetails struct {
	Tile string `json:"tile"`
	Row  int    `json:"row"`
	Col  int    `json:"col"`
}

type SnapshotDetails struct {
	Board          *board
	Hands          map[string][]tile
	Points         map[string]int
	Round          int
	RoundsUntilEnd int
	Variant        string
}

var (
	// list of all the unqiue tile paths
	uniquePaths = []string{
		A + F + B + C + D + E,
		B + E + C + F + D + A,
		A + F + B + E + C + D,
		A + B + C + E + D + F,
		A + D + B + F + C + E,
	}

	// map from paths to number of copies
	numCopiesByUniquePathsIndex = []int{6, 6, 14, 14, 14}

	// map from treature tile paths to (row, col) location
	initTreasureTiles = map[string][2]int{
		C + E + D: {0, 0},
		D + F + E: {0, 4},
		A + E + F: {4, 8},
		B + F + A: {8, 4},
		A + C + B: {8, 0},
		B + D + C: {4, 0},
		Special:   {4, 4},
	}

	// inital gems
	initGems = [][4]interface{}{
		{Amber, D, 0, 0},
		{Amber, E, 0, 4},
		{Amber, F, 4, 8},
		{Amber, A, 8, 4},
		{Amber, B, 8, 0},
		{Amber, C, 4, 0},
		{Emerald, Special, 4, 4},
		{Emerald, Special, 4, 4},
		{Emerald, Special, 4, 4},
		{Emerald, Special, 4, 4},
		{Emerald, Special, 4, 4},
		{Sapphire, Special, 4, 4},
	}

	// map from edges to (row, col) locations of every gateway
	initGateways = map[string][3][2]int{
		A + B: {{0, 1}, {0, 2}, {0, 3}},
		B + C: {{1, 5}, {2, 6}, {3, 7}},
		C + D: {{5, 7}, {6, 6}, {7, 5}},
		D + E: {{8, 3}, {8, 2}, {8, 1}},
		E + F: {{7, 0}, {6, 0}, {5, 0}},
		F + A: {{3, 0}, {2, 0}, {1, 0}},
	}

	// map from len of teams list to map of which teams own which gateways
	numTeamsToGatewayOwnership = map[int]map[string][]int{
		2: {
			A + B: {0},
			B + C: {1},
			C + D: {0},
			D + E: {1},
			E + F: {0},
			F + A: {1},
		},
		3: {
			A + B: {0},
			B + C: {0, 1},
			C + D: {2},
			D + E: {2, 0},
			E + F: {1},
			F + A: {1, 2},
		},
		4: {
			A + B: {0, 1},
			B + C: {1, 2},
			C + D: {0, 3},
			D + E: {3, 1},
			E + F: {2, 0},
			F + A: {2, 3},
		},
	}
)
