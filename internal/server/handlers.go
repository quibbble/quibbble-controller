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
	gs.mux.ServeHTTP(w, r)
}

func (gs *GameServer) connectHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" || gs.isConnected(name) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // allow origin checks are handled at the ingress
	})
	if err != nil {
		log.Println(err.Error())
		return
	}

	p := NewPlayer(name, conn, gs.actionCh)
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
	switch r.URL.Query().Get(FormatQuery) {
	case JSONFormat:
		snapshot, err := gs.GetSnapshotJSON()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		raw, _ := json.Marshal(snapshot)
		w.Header().Set("Content-Type", "application/json")
		w.Write(raw)
	case QGNFormat:
		snapshot, err := gs.GetSnapshotQGN()
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

type Activity struct {
	PlayerCount int       `json:"player_count"`
	LastUpdated time.Time `json:"last_updated"`
}

func (gs *GameServer) activeHandler(w http.ResponseWriter, r *http.Request) {
	raw, _ := json.Marshal(Activity{
		PlayerCount: len(gs.connected),
		LastUpdated: gs.lastUpdated,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(raw))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
