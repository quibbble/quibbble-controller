package carcassonne

import "github.com/quibbble/quibbble-controller/pkg/util"

// Action types
const (
	PlaceTileAction  = "place_tile"
	PlaceTokenAction = "place_token"
	RotateAction     = "rotate"
)

var (
	ActionToQGN = map[string]string{
		PlaceTileAction:  "i",
		PlaceTokenAction: "o",
		RotateAction:     "r",
	}
	QGNToAction = util.ReverseMap(ActionToQGN)

	sideToNotation = map[string]string{SideTop: "t", SideRight: "r", SideBottom: "b", SideLeft: "l"}
	notationToSide = util.ReverseMap(sideToNotation)

	farmSideToNotation = map[string]string{FarmSideTopA: "ta", FarmSideTopB: "tb", FarmSideRightA: "ra", FarmSideRightB: "rb", FarmSideBottomA: "ba", FarmSideBottomB: "bb", FarmSideLeftA: "la", FarmSideLeftB: "lb"}
	notationToFarmSide = util.ReverseMap(farmSideToNotation)

	structureToNotation = map[string]string{Road: "r", Farm: "f", City: "c", Cloister: "m", NilStructure: "n"}
	notationToStructure = util.ReverseMap(structureToNotation)

	tokenToNotation = map[string]string{Farmer: "f", Knight: "k", Thief: "t", Monk: "m"}
	notationToToken = util.ReverseMap(tokenToNotation)

	boolToNotation = map[bool]string{true: "t", false: "f"}
	notationToBool = map[string]bool{"t": true, "f": false}
)

// PlaceTileDetails is the action details for placing a tile
type PlaceTileDetails struct {
	// X and Y location where to place the tile
	X int `json:"x"`
	Y int `json:"y"`

	// Tile is the tile being placed
	Tile Tile `json:"tile"`
}

type Tile struct {
	Top             string `json:"top"`
	Right           string `json:"right"`
	Bottom          string `json:"bottom"`
	Left            string `json:"left"`
	Center          string `json:"center"`
	ConnectedCities bool   `json:"connected_cities" mapstructure:"connected_cities"`
	Banner          bool   `json:"banner"`
}

// PlaceTokenDetails is the action details for placing a token
type PlaceTokenDetails struct {
	// Pass set to pass placing a token
	Pass bool `json:"pass"`

	// X and Y location where to place the token
	X int `json:"X"`
	Y int `json:"y"`

	// Type is the type of token to place
	Type string `json:"type"`

	// Side is the side to place
	Side string `json:"side"`
}

// SnapshotDetails is the game data unique to Carcassonne
type SnapshotDetails struct {
	PlayTile       *tile            `json:"play_tile"`
	LastPlaced     map[string]*tile `json:"last_placed"`
	Board          []*tile          `json:"board"`
	BoardTokens    []*token         `json:"board_tokens"`
	Tokens         map[string]int   `json:"tokens"`
	Scores         map[string]int   `json:"scores"`
	TilesRemaining int              `json:"tiles_remaining"`
}

// startTile the tile at 0,0 at the start of the game
var startTile = newTile(City, Road, Farm, Road, NilStructure, false, false)

type tileAmounts struct {
	tile   *tile
	amount int
}

// tiles are all the tiles that will be placed
var tiles = []*tileAmounts{
	{tile: newTile(Farm, Farm, Farm, Farm, Cloister, false, false), amount: 4},
	{tile: newTile(Farm, Farm, Road, Farm, Cloister, false, false), amount: 2},
	{tile: newTile(City, City, City, City, NilStructure, true, true), amount: 1},
	{tile: newTile(City, City, Farm, City, NilStructure, true, false), amount: 3},
	{tile: newTile(City, City, Farm, City, NilStructure, true, true), amount: 1},
	{tile: newTile(City, City, Road, City, NilStructure, true, false), amount: 1},
	{tile: newTile(City, City, Road, City, NilStructure, true, true), amount: 2},
	{tile: newTile(City, Farm, Farm, City, NilStructure, true, false), amount: 3},
	{tile: newTile(City, Farm, Farm, City, NilStructure, true, true), amount: 2},
	{tile: newTile(City, Road, Road, City, NilStructure, true, false), amount: 3},
	{tile: newTile(City, Road, Road, City, NilStructure, true, true), amount: 2},
	{tile: newTile(Farm, City, Farm, City, NilStructure, true, false), amount: 1},
	{tile: newTile(Farm, City, Farm, City, NilStructure, true, true), amount: 2},
	{tile: newTile(City, Farm, Farm, City, NilStructure, false, false), amount: 2},
	{tile: newTile(City, Farm, City, Farm, NilStructure, false, false), amount: 3},
	{tile: newTile(City, Farm, Farm, Farm, NilStructure, false, false), amount: 5},
	{tile: newTile(City, Farm, Road, Road, NilStructure, false, false), amount: 3},
	{tile: newTile(City, Road, Road, Farm, NilStructure, false, false), amount: 3},
	{tile: newTile(City, Road, Road, Road, NilStructure, false, false), amount: 3},
	{tile: newTile(City, Road, Farm, Road, NilStructure, false, false), amount: 3},
	{tile: newTile(Road, Farm, Road, Farm, NilStructure, false, false), amount: 8},
	{tile: newTile(Farm, Farm, Road, Road, NilStructure, false, false), amount: 9},
	{tile: newTile(Farm, Road, Road, Road, NilStructure, false, false), amount: 4},
	{tile: newTile(Road, Road, Road, Road, NilStructure, false, false), amount: 1},
}
