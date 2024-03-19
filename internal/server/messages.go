package server

import (
	"encoding/json"
)

// Message Types
const (
	ErrorMessage      = "error"
	SnapshotMessage   = "snapshot"
	ConnectionMessage = "connection"
	ChatMessage       = "chat"
)

// Error Messages
var (
	ErrInvalidActionMessage     []byte
	ErrInvalidActionTypeMessage []byte
)

func init() {
	ErrInvalidActionMessage, _ = json.Marshal(Message{
		Type:    ErrorMessage,
		Details: "invalid action",
	})
	ErrInvalidActionTypeMessage, _ = json.Marshal(Message{
		Type:    ErrorMessage,
		Details: "invalid action type",
	})
}

type Message struct {
	Type    string      `json:"type"`
	Details interface{} `json:"details"`
}

func (gs *GameServer) sendSnapshotMessage(player *Player) {
	snapshot, _ := gs.game.GetSnapshotJSON()
	if player.team != nil {
		snapshot, _ = gs.game.GetSnapshotJSON(*player.team)
	}
	payload, _ := json.Marshal(Message{
		Type:    SnapshotMessage,
		Details: snapshot,
	})
	gs.sendMessage(player, payload)
}

func (gs *GameServer) sendSnapshotMessages() {
	for p := range gs.players {
		gs.sendSnapshotMessage(p)
	}
}

func (gs *GameServer) sendConnectionMessages() {
	players := make(map[string]*string)
	for player := range gs.players {
		players[player.uid] = player.team
	}
	for p := range gs.players {
		payload, _ := json.Marshal(Message{
			Type: ConnectionMessage,
			Details: struct {
				UID     string             `json:"uid"`
				Players map[string]*string `json:"players"`
			}{
				UID:     p.uid,
				Players: players,
			},
		})
		gs.sendMessage(p, payload)
	}
}

func (gs *GameServer) sendChatMessages(player *Player, message string) {
	payload, _ := json.Marshal(Message{
		Type: ChatMessage,
		Details: struct {
			UID     string  `json:"uid"`
			Team    *string `json:"team"`
			Message string  `json:"message"`
		}{
			UID:     player.uid,
			Team:    player.team,
			Message: message,
		},
	})
	for p := range gs.players {
		gs.sendMessage(p, payload)
	}
}

func (gs *GameServer) sendErrorMessage(player *Player, err error) {
	payload, _ := json.Marshal(Message{
		Type:    ErrorMessage,
		Details: err.Error(),
	})
	gs.sendMessage(player, payload)
}

func (gs *GameServer) sendMessage(player *Player, payload []byte) {
	select {
	case player.messageCh <- payload:
	default:
		delete(gs.players, player)
		go player.Close()
	}
}
