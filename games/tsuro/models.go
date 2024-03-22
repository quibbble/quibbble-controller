package tsuro

import "github.com/quibbble/quibbble-controller/pkg/util"

const BoardSize = 3

const (
	PlaceAction  = "place"
	RotateAction = "rotate"
)

var Actions = []string{PlaceAction, RotateAction}

const (
	ClassicVariant       = "classic"        // Standard Tsuro.
	LongestPathVariant   = "longest_path"   // Player with the longest path wins.
	MostCrossingsVariant = "most_crossings" // Player whose path crosses itself the most wins.
	OpenTilesVariant     = "open_tiles"     // Tiles are shared globally,
	SoloVariant          = "solo"           // Place tiles while keeping all tokens on the board,
)

var Variants = []string{ClassicVariant, LongestPathVariant, MostCrossingsVariant, OpenTilesVariant, SoloVariant}

var (
	ActionToQGN = map[string]string{
		PlaceAction: "p",
	}
	QGNToAction = util.ReverseMap(ActionToQGN)
)

type PlaceDetails struct {
	Row  int    `json:"row"`
	Col  int    `json:"col"`
	Tile string `json:"tile"`
}

type RotateDetails struct {
	Tile string `json:"tile"`
}

type SnapshotDetails struct {
	Variant        string             `json:"variant"`
	Board          [][]*tile          `json:"board"`
	TilesRemaining int                `json:"tiles_remaining"`
	Hands          map[string][]*tile `json:"hands"`
	Tokens         map[string]*token  `json:"tokens"`
	Dragon         string             `json:"dragon,omitempty"`
	Points         map[string]int     `json:"points,omitempty"`
}

// list of all the tiles that can be played
var tiles = []string{
	"ABCDEFGH", "AHBGCDEF", "AHBCDGEF", "AHBCDEFG", "AGBHCDEF",
	"ABCHDGEF", "ABCGDHEF", "AGBCDHEF", "ABCGDEFH", "AGBCDEFH",
	"ACBGDEFH", "ACBGDHEF", "ACBHDGEF", "ADBHCGEF", "ADBGCHEF",
	"ADBCEHFG", "ADBCEGFH", "AEBCDGFH", "AEBCDHFG", "AFBHCDEG",
	"AFBGCHDE", "AFBCDHEG", "AFBDCHEG", "AFBDCGEH", "AEBDCGFH",
	"ACBDEGFH", "AFBECHDG", "AFBECGDH", "AEBFCGDH", "ADBFCGEH",
	"ADBFCHEG", "ACBFDHEG", "ADBGCEFH", "AGBDCEFH", "ADBGCFEH",
}

// map of path to list of all other paths that cross the path
var crossing = map[string][]string{
	"AB": {}, "BA": {},
	"CD": {}, "DC": {},
	"EF": {}, "FE": {},
	"GH": {}, "HG": {},

	"AC": {"BD", "BE", "BF", "BG", "BH", "DB", "EB", "FB", "GB", "HB"},
	"CA": {"BD", "BE", "BF", "BG", "BH", "DB", "EB", "FB", "GB", "HB"},

	"AD": {"BE", "BF", "BG", "BH", "CE", "CF", "CG", "CH", "EB", "FB", "GB", "HB", "EC", "FC", "GC", "HC"},
	"DA": {"BE", "BF", "BG", "BH", "CE", "CF", "CG", "CH", "EB", "FB", "GB", "HB", "EC", "FC", "GC", "HC"},

	"AE": {"BF", "BG", "BH", "CF", "CG", "CH", "DF", "DG", "DH", "FB", "GB", "HB", "FC", "GC", "HC", "FD", "GD", "HD"},
	"EA": {"BF", "BG", "BH", "CF", "CG", "CH", "DF", "DG", "DH", "FB", "GB", "HB", "FC", "GC", "HC", "FD", "GD", "HD"},

	"AF": {"BG", "BH", "CG", "CH", "DG", "DH", "EG", "EH", "GB", "HB", "GC", "HC", "GD", "HD", "GE", "HE"},
	"FA": {"BG", "BH", "CG", "CH", "DG", "DH", "EG", "EH", "GB", "HB", "GC", "HC", "GD", "HD", "GE", "HE"},

	"AG": {"HB", "HC", "HD", "HE", "HF", "BH", "CH", "DH", "EH", "FH"},
	"GA": {"HB", "HC", "HD", "HE", "HF", "BH", "CH", "DH", "EH", "FH"},

	"AH": {}, "HA": {},

	"BC": {}, "CB": {},

	"BD": {"CA", "CE", "CF", "CG", "CH", "AC", "EC", "FC", "GC", "HC"},
	"DB": {"CA", "CE", "CF", "CG", "CH", "AC", "EC", "FC", "GC", "HC"},

	"BE": {"AC", "AD", "FC", "FD", "GC", "GD", "HC", "HD", "CA", "DA", "CF", "DF", "CG", "DG", "CH", "DH"},
	"EB": {"AC", "AD", "FC", "FD", "GC", "GD", "HC", "HD", "CA", "DA", "CF", "DF", "CG", "DG", "CH", "DH"},

	"BF": {"AE", "AC", "AD", "HE", "HC", "HD", "GE", "GC", "GD", "EA", "CA", "DA", "EH", "CH", "DH", "EG", "CG", "DG"},
	"FB": {"AE", "AC", "AD", "HE", "HC", "HD", "GE", "GC", "GD", "EA", "CA", "DA", "EH", "CH", "DH", "EG", "CG", "DG"},

	"BG": {"AC", "AD", "AE", "AF", "HC", "HD", "HE", "HF", "CA", "DA", "EA", "FA", "CH", "DH", "EH", "FH"},
	"GB": {"AC", "AD", "AE", "AF", "HC", "HD", "HE", "HF", "CA", "DA", "EA", "FA", "CH", "DH", "EH", "FH"},

	"BH": {"AC", "AD", "AE", "AF", "AG", "CA", "DA", "EA", "FA", "GA"},
	"HB": {"AC", "AD", "AE", "AF", "AG", "CA", "DA", "EA", "FA", "GA"},

	"CE": {"DA", "DB", "DF", "DG", "DH", "AD", "BD", "FD", "GD", "HD"},
	"EC": {"DA", "DB", "DF", "DG", "DH", "AD", "BD", "FD", "GD", "HD"},

	"CF": {"DA", "DB", "DG", "DH", "EA", "EB", "EG", "EH", "AD", "BD", "GD", "HD", "AE", "BE", "GE", "HE"},
	"FC": {"DA", "DB", "DG", "DH", "EA", "EB", "EG", "EH", "AD", "BD", "GD", "HD", "AE", "BE", "GE", "HE"},

	"CG": {"DH", "DA", "DB", "EH", "EA", "EB", "FH", "FA", "FB", "HD", "AD", "BD", "HE", "AE", "BE", "HF", "AF", "BF"},
	"GC": {"DH", "DA", "DB", "EH", "EA", "EB", "FH", "FA", "FB", "HD", "AD", "BD", "HE", "AE", "BE", "HF", "AF", "BF"},

	"CH": {"AD", "AE", "AF", "AG", "BD", "BE", "BF", "BG", "DA", "EA", "FA", "GA", "DB", "EB", "FB", "GB"},
	"HC": {"AD", "AE", "AF", "AG", "BD", "BE", "BF", "BG", "DA", "EA", "FA", "GA", "DB", "EB", "FB", "GB"},

	"DE": {}, "ED": {},

	"DF": {"EA", "EB", "EC", "EG", "EH", "AE", "BE", "CE", "GE", "HE"},
	"FD": {"EA", "EB", "EC", "EG", "EH", "AE", "BE", "CE", "GE", "HE"},

	"DG": {"EA", "EB", "EC", "EH", "FA", "FB", "FC", "FH", "AE", "AF", "BE", "BF", "CE", "CF", "HE", "HF"},
	"GD": {"EA", "EB", "EC", "EH", "FA", "FB", "FC", "FH", "AE", "AF", "BE", "BF", "CE", "CF", "HE", "HF"},

	"DH": {"GA", "GB", "GC", "FA", "FB", "FC", "EA", "EB", "EC", "AG", "BG", "CG", "AF", "BF", "CF", "AE", "BE", "CE"},
	"HD": {"GA", "GB", "GC", "FA", "FB", "FC", "EA", "EB", "EC", "AG", "BG", "CG", "AF", "BF", "CF", "AE", "BE", "CE"},

	"EG": {"FA", "FB", "FC", "FD", "FH", "AF", "BF", "CF", "DF", "HF"},
	"GE": {"FA", "FB", "FC", "FD", "FH", "AF", "BF", "CF", "DF", "HF"},

	"EH": {"FA", "FB", "FC", "FD", "GA", "GB", "GC", "GD", "AF", "BF", "CF", "DF", "AG", "BG", "CG", "DG"},
	"HE": {"FA", "FB", "FC", "FD", "GA", "GB", "GC", "GD", "AF", "BF", "CF", "DF", "AG", "BG", "CG", "DG"},

	"FG": {}, "GF": {},

	"FH": {"GA", "GB", "GC", "GD", "GE", "AG", "BG", "CG", "DG", "EG"},
	"HF": {"GA", "GB", "GC", "GD", "GE", "AG", "BG", "CG", "DG", "EG"},
}
