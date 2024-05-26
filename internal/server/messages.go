package server

import (
	"encoding/json"
)

// Message Types
const (
	PongMessage       = "pong"
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
	snapshot, _ := gs.GetSnapshotJSON()
	if team := gs.team(player.name); team != nil {
		snapshot, _ = gs.GetSnapshotJSON(*team)
	}
	payload, _ := json.Marshal(Message{
		Type:    SnapshotMessage,
		Details: snapshot,
	})
	gs.sendMessage(player, payload)
}

func (gs *GameServer) sendSnapshotMessages() {
	for p := range gs.connected {
		gs.sendSnapshotMessage(p)
	}
}

func (gs *GameServer) sendConnectionMessages() {
	connected := make(map[string]*string)
	for player := range gs.connected {
		connected[player.name] = player.team
	}
	for p := range gs.connected {
		payload, _ := json.Marshal(Message{
			Type:    ConnectionMessage,
			Details: connected,
		})
		gs.sendMessage(p, payload)
	}
}

func (gs *GameServer) sendChatMessages(player *Player, message string) {
	payload, _ := json.Marshal(Message{
		Type: ChatMessage,
		Details: struct {
			Name    string `json:"name"`
			Message string `json:"message"`
		}{
			Name:    player.name,
			Message: message,
		},
	})
	for p := range gs.connected {
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

func (gs *GameServer) sendPongMessage(player *Player) {
	payload, _ := json.Marshal(Message{
		Type: PongMessage,
	})
	gs.sendMessage(player, payload)
}

func (gs *GameServer) sendMessage(player *Player, payload []byte) {
	select {
	case player.messageCh <- payload:
	default:
		delete(gs.connected, player)
		go player.Close()
	}
}
