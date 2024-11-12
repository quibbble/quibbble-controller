package quill

import (
	"fmt"

	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

var (
	actionToNotation = map[string]string{
		PlayCardAction:   "p",
		SackCardAction:   "s",
		AttackUnitAction: "a",
		MoveUnitAction:   "m",
		EndTurnAction:    "e",
	}
	notationToAction = reverseMap(actionToNotation)
)

func (s *PlayCardActionDetails) encodeBGN() []string {
	str := make([]string, 0)
	for _, target := range s.Targets {
		str = append(str, string(target))
	}
	return append([]string{string(s.Card)}, str...)
}

func decodePlayCard(notation []string) (*PlayCardActionDetails, error) {
	if len(notation) < 1 {
		return nil, fmt.Errorf("invalid play card notation")
	}
	targets := make([]uuid.UUID, 0)
	for i := 1; i < len(notation); i++ {
		targets = append(targets, uuid.UUID(notation[i]))
	}
	return &PlayCardActionDetails{
		Card:    uuid.UUID(notation[0]),
		Targets: targets,
	}, nil
}

func (s *SackCardActionDetails) encodeBGN() []string {
	return append([]string{string(s.Card)}, s.Option)
}

func decodeSackCard(notation []string) (*SackCardActionDetails, error) {
	if len(notation) != 2 {
		return nil, fmt.Errorf("invalid sack card notation")
	}
	return &SackCardActionDetails{
		Card:   uuid.UUID(notation[0]),
		Option: notation[1],
	}, nil
}

func (s *AttackUnitActionDetails) encodeBGN() []string {
	return append([]string{string(s.Attacker)}, string(s.Defender))
}

func decodeAttackUnit(notation []string) (*AttackUnitActionDetails, error) {
	if len(notation) != 2 {
		return nil, fmt.Errorf("invalid attack unit notation")
	}
	return &AttackUnitActionDetails{
		Attacker: uuid.UUID(notation[0]),
		Defender: uuid.UUID(notation[1]),
	}, nil
}

func (s *MoveUnitActionDetails) encodeBGN() []string {
	return []string{string(s.Unit), string(s.Tile)}
}

func decodeMoveUnit(notation []string) (*MoveUnitActionDetails, error) {
	if len(notation) != 2 {
		return nil, fmt.Errorf("invalid move unit notation")
	}
	return &MoveUnitActionDetails{
		Unit: uuid.UUID(notation[0]),
		Tile: uuid.UUID(notation[1]),
	}, nil
}

func (s *EndTurnActionDetails) encodeBGN() []string {
	return []string{}
}

func decodeEndTurn(notation []string) (*EndTurnActionDetails, error) {
	if len(notation) != 2 {
		return nil, fmt.Errorf("invalid move unit notation")
	}
	return &EndTurnActionDetails{}, nil
}
