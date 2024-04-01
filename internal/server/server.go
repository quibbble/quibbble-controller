package server

import (
	"fmt"
	"log"
	"runtime/debug"
	"slices"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type GameServer struct {
	// createdAt and updatedAt all log time information about the server.
	lastUpdated time.Time

	// mux routes the various endpoints to the appropriate handler.
	mux *chi.Mux

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

	// allowedOrigins are the list of locations that may connect to the server
	allowedOrigins []string
}

func NewGameServer(game qg.Game, completeFn func(qg.Game), allowedOrigins []string) *GameServer {
	gs := &GameServer{
		lastUpdated:    time.Now(),
		mux:            chi.NewRouter(),
		game:           game,
		players:        make(map[*Player]struct{}),
		joinCh:         make(chan *Player),
		leaveCh:        make(chan *Player),
		actionCh:       make(chan *Action),
		completeFn:     completeFn,
		allowedOrigins: allowedOrigins,
	}
	go gs.Start()
	gs.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
	}))
	gs.mux.Get("/connect", gs.connectHandler)
	gs.mux.Get("/snapshot", gs.snapshotHandler)
	gs.mux.Get("/active", gs.activeHandler)
	gs.mux.Get("/health", healthHandler)
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
