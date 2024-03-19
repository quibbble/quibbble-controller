package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

func (gs *GameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gs.serveMux.ServeHTTP(w, r)
}

func (gs *GameServer) connectHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	p := NewPlayer(conn, gs.actionCh)
	gs.joinCh <- p

	ctx := context.Background()
	go func() {
		if err := p.ReadPump(ctx); err != nil {
			log.Println(err.Error())
		}
		gs.leaveCh <- p
		p.conn.CloseNow()
	}()
	go func() {
		if err := p.WritePump(ctx); err != nil {
			log.Println(err.Error())
		}
	}()
}

func (gs *GameServer) snapshotHandler(w http.ResponseWriter, r *http.Request) {
	const (
		FormatQuery = "format"
		JSONFormat  = "json"
		QGNFormat   = "qgn"
	)
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	switch r.URL.Query().Get("format") {
	case JSONFormat:
		snapshot, err := gs.game.GetSnapshotJSON()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		raw, _ := json.Marshal(snapshot)
		w.Header().Set("Content-Type", "application/json")
		w.Write(raw)
	case QGNFormat:
		snapshot, err := gs.game.GetSnapshotQGN()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/qgn")
		w.Write([]byte(snapshot.String()))
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

type Active struct {
	Players     int       `json:"players"`
	LastUpdated time.Time `json:"last_updated"`
}

func (gs *GameServer) activeHandler(w http.ResponseWriter, r *http.Request) {
	raw, _ := json.Marshal(Active{
		Players:     len(gs.players),
		LastUpdated: gs.lastUpdated,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(raw))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
