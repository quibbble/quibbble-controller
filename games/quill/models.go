package quill

import (
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

// Action types
const (
	// Targeting actions
	TargetAction      = "Target"
	NextTargetsAction = "NextTargets"

	// Normal actions
	PlayCardAction   = "PlayCard"
	SackCardAction   = "SackCard"
	AttackUnitAction = "AttackUnit"
	MoveUnitAction   = "MoveUnit"
	EndTurnAction    = "EndTurn"
)

const (
	DecksTag = "Decks"
)

type QuillInformation struct {
	Cards []st.ICard
}

type QuillSnapshotData struct {
	Board      [st.Cols][st.Rows]*st.Tile
	PlayRange  map[string][]int
	UUIDToTeam map[uuid.UUID]string
	Hand       map[string][]st.ICard
	Deck       map[string]int
	Mana       map[string]*st.Mana
	Sacked     map[string]bool
	Targets    []uuid.UUID
}

type NextTargetsActionDetails struct {
	Targets []uuid.UUID
}

type PlayCardActionDetails struct {
	Card    uuid.UUID
	Targets []uuid.UUID

	// DO NOT SET - QUILL INTERNAL USE ONLY
	PlayCard st.ICard
}

type SackCardActionDetails struct {
	Card   uuid.UUID
	Option string
}

type AttackUnitActionDetails struct {
	Attacker, Defender uuid.UUID

	// DO NOT SET - QUILL INTERNAL USE ONLY
	AttackerCard, DefenderCard st.ICard
}

type MoveUnitActionDetails struct {
	Unit, Tile uuid.UUID

	// DO NOT SET - QUILL INTERNAL USE ONLY
	UnitCard st.ICard
	TileXY   []int
}

type EndTurnActionDetails struct{}
