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
		snapshot, _ = gs.GetSnapshotJSON(*player.team)
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
	teams := make(map[string]*string)
	usernames := make(map[string]string)
	for player := range gs.connected {
		teams[player.uid] = player.team
		usernames[player.uid] = player.username
	}
	for p := range gs.connected {
		payload, _ := json.Marshal(Message{
			Type: ConnectionMessage,
			Details: struct {
				UID       string             `json:"uid"`
				Teams     map[string]*string `json:"teams"`
				Usernames map[string]string  `json:"usernames"`
			}{
				UID:       p.uid,
				Teams:     teams,
				Usernames: usernames,
			},
		})
		gs.sendMessage(p, payload)
	}
}

func (gs *GameServer) sendChatMessages(player *Player, message string) {
	payload, _ := json.Marshal(Message{
		Type: ChatMessage,
		Details: struct {
			UID      string  `json:"uid"`
			Username string  `json:"username"`
			Team     *string `json:"team"`
			Message  string  `json:"message"`
		}{
			UID:      player.uid,
			Username: player.username,
			Team:     player.team,
			Message:  message,
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

func (gs *GameServer) sendMessage(player *Player, payload []byte) {
	select {
	case player.messageCh <- payload:
	default:
		delete(gs.connected, player)
		go player.Close()
	}
}
