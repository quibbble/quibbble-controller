package stratego

import "github.com/quibbble/quibbble-controller/pkg/util"

const (
	BoardSize            = 10
	QuickBattleBoardSize = 8
)

const (
	SwitchAction = "switch"
	MoveAction   = "move"
)

var Actions = []string{SwitchAction, MoveAction}

var (
	ActionToQGN = map[string]string{
		SwitchAction: "s",
		MoveAction:   "m",
	}
	QGNToAction = util.ReverseMap(ActionToQGN)
)

const (
	ClassicVariant     = "classic"      // normal Stratego
	QuickBattleVariant = "quick_battle" // 8x8 quick play Stratego
)

var Variants = []string{ClassicVariant, QuickBattleVariant}

type SwitchDetails struct {
	UnitARow int `json:"unita_row" mapstructure:"unita_row"`
	UnitACol int `json:"unita_col" mapstructure:"unita_col"`

	UnitBRow int `json:"unitb_row" mapstructure:"unitb_row"`
	UnitBCol int `json:"unitb_col" mapstructure:"unitb_col"`
}

type MoveDetails struct {
	UnitRow int `json:"unit_row" mapstructure:"unit_row"`
	UnitCol int `json:"unit_col" mapstructure:"unit_col"`

	TileRow int `json:"tile_row" mapstructure:"tile_row"`
	TileCol int `json:"tile_col" mapstructure:"tile_col"`
}

type Battle struct {
	MoveDetails
	AttackingUnit Unit   `json:"attacking_unit" mapstructure:"attacking_unit"`
	AttackedUnit  Unit   `json:"attacked_unit" mapstructure:"attacked_unit"`
	WinningTeam   string `json:"winning_team" mapstructure:"winning_team"`
}

type SnapshotDetails struct {
	Board       [][]Unit `json:"board"`
	Battle      *Battle  `json:"battle"`
	JustBattled bool     `json:"just_battled"`
	Started     bool     `json:"started"`
	Variant     string   `json:"variant"`
}
