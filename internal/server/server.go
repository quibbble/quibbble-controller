package server

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"slices"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

type GameServer struct {
	// createdAt and updatedAt all log time information about the server.
	lastUpdated time.Time

	// mux routes the various endpoints to the appropriate handler.
	mux http.ServeMux

	// game is the instance of the game being played.
	game qg.Game

	// id of the game server
	id string

	// kind of game
	kind string

	// connected represents all players currently connected to the server.
	connected map[*Player]struct{}

	// joinCh and leaveCh adds/remove a player from the server.
	joinCh, leaveCh chan *Player

	// actionCh sends actions to the server to be processed.
	actionCh chan *Action

	// completeFn is called on game end.
	completeFn func(qg.Game)
}

func NewGameServer(game qg.Game, id, kind string, completeFn func(qg.Game)) *GameServer {
	gs := &GameServer{
		lastUpdated: time.Now(),
		game:        game,
		id:          id,
		kind:        kind,
		connected:   make(map[*Player]struct{}),
		joinCh:      make(chan *Player),
		leaveCh:     make(chan *Player),
		actionCh:    make(chan *Action),
		completeFn:  completeFn,
	}
	go gs.Start()
	// these will be prefixed by /game/{key}/{id} when being called through nginx
	gs.mux.HandleFunc("GET /", gs.connectHandler)
	gs.mux.HandleFunc("GET /snapshot", gs.snapshotHandler)
	gs.mux.HandleFunc("GET /activity", gs.activeHandler)
	gs.mux.HandleFunc("GET /health", healthHandler)
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
			gs.sendConnectionMessages()
			gs.sendSnapshotMessage(p)
		case p := <-gs.leaveCh:
			delete(gs.connected, p)
			go p.Close()
			gs.sendConnectionMessages()
		case a := <-gs.actionCh:
			switch a.Type {
			case Ping:
				gs.sendPongMessage(a.Player)
				continue
			case Chat:
				message, ok := a.Details.(string)
				if !ok {
					gs.sendErrorMessage(a.Player, fmt.Errorf("invalid chat action"))
					continue
				}
				gs.sendChatMessages(a.Player, message)
				continue
			case Join:
				team, ok := a.Details.(string)
				if !ok {
					gs.sendErrorMessage(a.Player, fmt.Errorf("invalid join action"))
					continue
				}
				snapshot, err := gs.game.GetSnapshotJSON()
				if err != nil {
					gs.sendErrorMessage(a.Player, fmt.Errorf("internal game failure"))
					continue
				}
				if !slices.Contains(snapshot.Teams, team) {
					gs.sendErrorMessage(a.Player, fmt.Errorf("invalid team"))
					continue
				}
				a.Player.team = &team
				gs.sendConnectionMessages()
			default:
				team := a.Player.team
				if team == nil {
					gs.sendErrorMessage(a.Player, fmt.Errorf("not part of a team"))
					continue
				}
				a.Action.Team = *team
				if err := gs.game.Do(a.Action); err != nil {
					gs.sendErrorMessage(a.Player, err)
					continue
				}
				gs.sendSnapshotMessages()

				if snapshot, err := gs.GetSnapshotJSON(); err == nil && len(snapshot.Winners) > 0 {
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
	snapshot.Kind = gs.kind
	return snapshot, nil
}

func (gs *GameServer) GetSnapshotQGN() (*qgn.Snapshot, error) {
	snapshot, err := gs.game.GetSnapshotQGN()
	if err != nil {
		return nil, err
	}
	// add missing server specific tags
	snapshot.Tags[qgn.IDTag] = gs.id
	snapshot.Tags[qgn.KindTag] = gs.kind
	return snapshot, nil
}

func (gs *GameServer) team(name string) *string {
	for player := range gs.connected {
		if player.name == name {
			return player.team
		}
	}
	return nil
}

func (gs *GameServer) isConnected(name string) bool {
	for player := range gs.connected {
		if player.name == name {
			return true
		}
	}
	return false
}
