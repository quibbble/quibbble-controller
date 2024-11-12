package quill

import (
	"fmt"
	"strings"

	"github.com/quibbble/quibbble-controller/games/quill/internal/game"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

const key = "Quill"

type Builder struct{}

func (b Builder) Create(snapshot *qgn.Snapshot) (qg.Game, error) {
	if snapshot.Tags[qgn.KeyTag] != key {
		return nil, fmt.Errorf("key tag does not match game key")
	}
	teams, err := snapshot.Tags.Teams()
	if err != nil {
		return nil, err
	}
	seed, err := snapshot.Tags.Seed()
	if err != nil {
		return nil, err
	}
	decksStr, ok := snapshot.Tags[DecksTag]
	if !ok {
		return nil, fmt.Errorf("missing decks tag")
	}
	decksList := strings.Split(decksStr, " : ")
	if len(decksList) != 2 {
		return nil, fmt.Errorf("must have two decks")
	}
	decks := [][]string{strings.Split(decksList[0], ", "), strings.Split(decksList[1], ", ")}
	game, err := NewQuill(seed, teams, decks)
	if err != nil {
		return nil, err
	}
	for _, action := range snapshot.Actions {
		if action.Index >= len(teams) {
			return nil, fmt.Errorf("team index %d out of range", action.Index)
		}
		team := teams[action.Index]
		actionType := notationToAction[string(action.Key)]
		if actionType == "" {
			return nil, fmt.Errorf("invalid action key %s", string(action.Key))
		}
		var details interface{}
		switch actionType {
		case PlayCardAction:
			result, err := decodePlayCard(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case SackCardAction:
			result, err := decodeSackCard(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case AttackUnitAction:
			result, err := decodeAttackUnit(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case MoveUnitAction:
			result, err := decodeMoveUnit(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case EndTurnAction:
			result, err := decodeEndTurn(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		}
		if err := game.Do(&qg.Action{
			Team:    team,
			Type:    actionType,
			Details: details,
		}); err != nil {
			return nil, err
		}
	}
	return game, nil
}

func (b Builder) GetInformation() *qg.Information {
	ids, err := parse.AllCards()
	if err != nil {
		return nil
	}
	cards := make([]st.ICard, 0)
	for _, id := range ids {
		card, err := game.NewDummyCard(id)
		if err != nil {
			if err == parse.ErrNotEnabled {
				continue
			}
			return nil
		}
		cards = append(cards, card)
	}
	return &qg.Information{
		Key: key,
		Min: minTeams,
		Max: maxTeams,
		Details: &QuillInformation{
			Cards: cards,
		},
	}
}
