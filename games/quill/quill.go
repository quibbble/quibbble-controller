package quill

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/quibbble/quibbble-controller/games/quill/internal/game"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const (
	minTeams = 2
	maxTeams = 2
)

type Quill struct {
	seed    int64
	decks   [][]string
	state   *game.Game
	history []*qg.Action

	teams      []string
	teamToUUID map[string]uuid.UUID
	uuidToTeam map[uuid.UUID]string
	targets    []uuid.UUID
}

func NewQuill(seed int64, teams []string, decks [][]string) (*Quill, error) {
	if len(teams) < minTeams {
		return nil, fmt.Errorf("at least %d teams required to create a game of %s", minTeams, key)
	} else if len(teams) > maxTeams {
		return nil, fmt.Errorf("at most %d teams allowed to create a game of %s", maxTeams, key)
	} else if duplicates(teams) {
		return nil, fmt.Errorf("duplicate teams found")
	}
	teamToUUID := make(map[string]uuid.UUID)
	uuidToTeam := make(map[uuid.UUID]string)
	gen := uuid.NewGen(rand.New(rand.NewSource(seed)))
	for _, team := range teams {
		playerUUID := gen.New('P')
		teamToUUID[team] = playerUUID
		uuidToTeam[playerUUID] = team
	}
	state, err := game.NewGame(seed, teamToUUID[teams[0]], teamToUUID[teams[1]], decks[0], decks[1])
	if err != nil {
		return nil, err
	}
	targets, err := state.GetNextTargets(state.GetTurn())
	if err != nil {
		return nil, err
	}
	return &Quill{
		seed:       seed,
		decks:      decks,
		state:      state,
		history:    make([]*qg.Action, 0),
		targets:    targets,
		teams:      teams,
		teamToUUID: teamToUUID,
		uuidToTeam: uuidToTeam,
	}, nil
}

func (q *Quill) Do(action *qg.Action) error {
	if q.state.Winner != nil {
		return fmt.Errorf("game already over")
	}
	team, ok := q.teamToUUID[action.Team]
	if !ok {
		return fmt.Errorf("team not found")
	}
	switch action.Type {
	case NextTargetsAction:
		var details NextTargetsActionDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		targets, err := q.state.GetNextTargets(team, details.Targets...)
		if err != nil {
			return err
		}
		q.targets = targets
	case PlayCardAction:
		var details PlayCardActionDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		details.PlayCard = q.state.GetCard(details.Card)
		action.Details = details
		if err := q.state.PlayCard(team, details.Card, details.Targets...); err != nil {
			return err
		}
		q.history = append(q.history, action)
		targets, err := q.state.GetNextTargets(team)
		if err != nil {
			return err
		}
		q.targets = targets
	case SackCardAction:
		var details SackCardActionDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		if err := q.state.SackCard(team, details.Card, details.Option); err != nil {
			return err
		}
		q.history = append(q.history, action)
		targets, err := q.state.GetNextTargets(team)
		if err != nil {
			return err
		}
		q.targets = targets
	case AttackUnitAction:
		var details AttackUnitActionDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		details.AttackerCard = q.state.GetCard(details.Attacker)
		details.DefenderCard = q.state.GetCard(details.Defender)
		action.Details = details
		if err := q.state.AttackUnit(team, details.Attacker, details.Defender); err != nil {
			return err
		}
		q.history = append(q.history, action)
		targets, err := q.state.GetNextTargets(team)
		if err != nil {
			return err
		}
		q.targets = targets
	case MoveUnitAction:
		var details MoveUnitActionDetails
		if err := mapstructure.Decode(action.Details, &details); err != nil {
			return err
		}
		details.UnitCard = q.state.GetCard(details.Unit)
		x, y, _ := q.state.Board.GetTileXY(details.Tile)
		details.TileXY = []int{x, y}
		action.Details = details
		if err := q.state.MoveUnit(team, details.Unit, details.Tile); err != nil {
			return err
		}
		q.history = append(q.history, action)
		targets, err := q.state.GetNextTargets(team)
		if err != nil {
			return err
		}
		q.targets = targets
	case EndTurnAction:
		if err := q.state.EndTurn(q.teamToUUID[action.Team]); err != nil {
			return err
		}
		q.history = append(q.history, action)
		targets, err := q.state.GetNextTargets(q.state.GetTurn())
		if err != nil {
			return err
		}
		q.targets = targets
	default:
		return fmt.Errorf("cannot process action type %s", action.Type)
	}
	return nil
}

func (q *Quill) GetSnapshotJSON(team ...string) (*qg.Snapshot, error) {
	if len(team) > 1 {
		return nil, fmt.Errorf("get snapshot requires zero or one team")
	}
	winners := make([]string, 0)
	if q.state.Winner != nil {
		winners = append(winners, q.uuidToTeam[*q.state.Winner])
	}
	hand := make(map[string][]st.ICard)
	for id, h := range q.state.Hand {
		cards := h.GetItems()
		if len(team) == 1 && q.uuidToTeam[id] != team[0] {
			empty := make([]st.ICard, 0)
			for i := 0; i < len(cards); i++ {
				empty = append(empty, cd.NewEmptyCard(id))
			}
			cards = empty
		}
		hand[q.uuidToTeam[id]] = cards
	}
	playRange := make(map[string][]int)
	for _, team := range q.teams {
		min, max := q.state.Board.GetPlayableRowRange(q.teamToUUID[team])
		playRange[team] = []int{min, max}
	}
	deck := make(map[string]int)
	for id, d := range q.state.Deck {
		deck[q.uuidToTeam[id]] = d.GetSize()
	}
	mana := make(map[string]*st.Mana)
	for id, m := range q.state.Mana {
		mana[q.uuidToTeam[id]] = m
	}
	sacked := make(map[string]bool)
	for id, s := range q.state.Sacked {
		sacked[q.uuidToTeam[id]] = s
	}
	targets := make([]uuid.UUID, 0)
	if len(team) == 1 && q.state.GetTurn() == q.teamToUUID[team[0]] {
		targets = q.targets
	}
	return &qg.Snapshot{
		Turn:    q.uuidToTeam[q.state.GetTurn()],
		Teams:   q.teams,
		Winners: winners,
		Details: QuillSnapshotData{
			Board:      q.state.Board.XYs,
			PlayRange:  playRange,
			UUIDToTeam: q.uuidToTeam,
			Hand:       hand,
			Deck:       deck,
			Mana:       mana,
			Sacked:     sacked,
			Targets:    targets,
		},
		// Actions: , // TODO calculate all valid actions and add here
		History: q.history,
		Message: fmt.Sprintf("%s must complete their turn", q.uuidToTeam[q.state.GetTurn()]),
	}, nil
}

func (q *Quill) GetSnapshotQGN() (*qgn.Snapshot, error) {
	tags := map[string]string{
		qgn.KindTag:  key,
		qgn.TeamsTag: strings.Join(q.teams, ", "),
		qgn.SeedTag:  fmt.Sprintf("%d", q.seed),
		DecksTag: strings.Join([]string{
			strings.Join(q.decks[0], ", "),
			strings.Join(q.decks[1], ", "),
		}, " : "),
	}
	actions := make([]qgn.Action, 0)
	for _, action := range q.history {
		qgnAction := qgn.Action{
			Index: indexOf(q.teams, action.Team),
			Key:   actionToNotation[action.Type],
		}
		switch action.Type {
		case PlayCardAction:
			var details PlayCardActionDetails
			_ = mapstructure.Decode(action.Details, &details)
			qgnAction.Details = details.encodeBGN()
		case SackCardAction:
			var details SackCardActionDetails
			_ = mapstructure.Decode(action.Details, &details)
			qgnAction.Details = details.encodeBGN()
		case AttackUnitAction:
			var details AttackUnitActionDetails
			_ = mapstructure.Decode(action.Details, &details)
			qgnAction.Details = details.encodeBGN()
		case MoveUnitAction:
			var details MoveUnitActionDetails
			_ = mapstructure.Decode(action.Details, &details)
			qgnAction.Details = details.encodeBGN()
		case EndTurnAction:
			var details EndTurnActionDetails
			_ = mapstructure.Decode(action.Details, &details)
			qgnAction.Details = details.encodeBGN()
		}
		actions = append(actions, qgnAction)
	}
	return &qgn.Snapshot{
		Tags:    tags,
		Actions: actions,
	}, nil
}
