package server

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

type GameServer struct {
	// createdAt and updatedAt all log time information about the server.
	lastUpdated time.Time

	// mux routes the various endpoints to the appropriate handler.
	mux *http.ServeMux

	// game is the instance of the game being played.
	game qg.Game

	// id of the game server
	id string

	// typ is one of ai, multiplayer, or local.
	typ string

	// players is a map from team to list of players allowed to join that team.
	players map[string][]string

	// connected represents all players currently connected to the server.
	connected map[*Player]struct{}

	// joinCh and leaveCh adds/remove a player from the server.
	joinCh, leaveCh chan *Player

	// actionCh sends actions to the server to be processed.
	actionCh chan *Action

	// completeFn is called on game end.
	completeFn func(qg.Game)
}

func NewGameServer(game qg.Game, id, typ string, players map[string][]string, completeFn func(qg.Game), authenticate func(http.Handler) http.Handler) *GameServer {
	gs := &GameServer{
		lastUpdated: time.Now(),
		game:        game,
		id:          id,
		typ:         typ,
		players:     players,
		connected:   make(map[*Player]struct{}),
		joinCh:      make(chan *Player),
		leaveCh:     make(chan *Player),
		actionCh:    make(chan *Action),
		completeFn:  completeFn,
	}
	go gs.Start()
	gs.mux.Handle("/connect", authenticate(http.HandlerFunc(gs.connectHandler)))
	gs.mux.HandleFunc("/snapshot", gs.snapshotHandler)
	gs.mux.HandleFunc("/active", gs.activeHandler)
	gs.mux.HandleFunc("/health", healthHandler)
	return gs
}

func (gs *GameServer) Start() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(string(debug.Stack()))
		}
	}()

	for {
		select {
		case p := <-gs.joinCh:
			gs.connected[p] = struct{}{}
			if gs.typ == qgn.AIType || gs.typ == qgn.MultiplayerType {
				p.team = team(gs.players, p.uid)
			}
			gs.sendConnectionMessages()
			gs.sendSnapshotMessage(p)
		case p := <-gs.leaveCh:
			delete(gs.connected, p)
			go p.Close()
			gs.sendConnectionMessages()
		case a := <-gs.actionCh:
			switch a.Type {
			case Join:
				if gs.typ == qgn.AIType || gs.typ == qgn.MultiplayerType {
					gs.sendErrorMessage(a.Player, fmt.Errorf("join action disabled for this game type"))
					continue
				}
				snapshot, err := gs.game.GetSnapshotJSON()
				if err != nil {
					gs.sendErrorMessage(a.Player, err)
					continue
				}
				team, ok := a.Details.(string)
				if !ok || !slices.Contains(snapshot.Teams, team) {
					gs.sendMessage(a.Player, ErrInvalidActionMessage)
					continue
				}
				a.Player.team = &team
				gs.sendConnectionMessages()
				gs.sendSnapshotMessage(a.Player)
				continue
			case Chat:
				message, ok := a.Details.(string)
				if !ok {
					gs.sendErrorMessage(a.Player, fmt.Errorf("invalid chat action"))
					continue
				}
				gs.sendChatMessages(a.Player, message)
				continue
			default:
				if a.Player.team == nil {
					gs.sendErrorMessage(a.Player, fmt.Errorf("not part of a team"))
					continue
				}
				a.Action.Team = *a.Player.team
				if err := gs.game.Do(a.Action); err != nil {
					gs.sendErrorMessage(a.Player, err)
					continue
				}
				gs.sendSnapshotMessages()

				if snapshot, err := gs.game.GetSnapshotJSON(); err == nil && len(snapshot.Winners) > 0 {
					gs.completeFn(gs.game)
				}
			}
		}
		gs.lastUpdated = time.Now()
	}
}

func (gs *GameServer) GetSnapshotJSON(team ...string) (*qg.Snapshot, error) {
	snapshot, err := gs.game.GetSnapshotJSON(team...)
	if err != nil {
		return nil, err
	}
	// add missing server specific data
	snapshot.Type = gs.typ
	return snapshot, nil
}

func (gs *GameServer) GetSnapshotQGN() (*qgn.Snapshot, error) {
	snapshot, err := gs.game.GetSnapshotQGN()
	if err != nil {
		return nil, err
	}
	// add missing server specific tags
	snapshot.Tags[qgn.IDTag] = gs.id
	snapshot.Tags[qgn.TypeTag] = gs.typ
	for team, players := range gs.players {
		tag := fmt.Sprintf("%s_%s", team, qgn.PlayersTagSuffix)
		snapshot.Tags[tag] = strings.Join(players, ", ")
	}
	return snapshot, nil
}

func team(players map[string][]string, uid string) *string {
	var team *string
	for t, players := range players {
		if slices.Contains(players, uid) {
			team = &t
			break
		}
	}
	return team
}
