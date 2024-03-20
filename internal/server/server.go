package server

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"slices"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type GameServer struct {
	// createdAt and updatedAt all log time information about the server.
	lastUpdated time.Time

	// serveMux routes the various endpoints to the appropriate handler.
	serveMux http.ServeMux

	// game is the instance of the game being played.
	game qg.Game

	// players is a map from player to team
	players map[*Player]struct{}

	// joinCh and leaveCh adds/remove a player from the server.
	joinCh, leaveCh chan *Player

	// actionCh sends actions to the server to be processed.
	actionCh chan *Action

	// completeFn is called on game end.
	completeFn func(qg.Game)
}

func NewGameServer(game qg.Game, completeFn func(qg.Game)) *GameServer {
	gs := &GameServer{
		lastUpdated: time.Now(),
		game:        game,
		players:     make(map[*Player]struct{}),
		joinCh:      make(chan *Player),
		leaveCh:     make(chan *Player),
		actionCh:    make(chan *Action),
		completeFn:  completeFn,
	}
	go gs.Start()
	gs.serveMux.HandleFunc("/connect", gs.connectHandler)
	gs.serveMux.HandleFunc("/snapshot", gs.snapshotHandler)
	gs.serveMux.HandleFunc("/active", gs.activeHandler)
	gs.serveMux.HandleFunc("/health", healthHandler)
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
			gs.players[p] = struct{}{}
			gs.sendConnectionMessages()
			gs.sendSnapshotMessage(p)
		case p := <-gs.leaveCh:
			delete(gs.players, p)
			go p.Close()
			gs.sendConnectionMessages()
		case a := <-gs.actionCh:
			switch a.Type {
			case Join:
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
