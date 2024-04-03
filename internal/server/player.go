package server

import (
	"context"
	"net"
	"sync"
	"time"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	// playerMessageBuffer controls the max number
	// of messages that can be queued for a player
	// before it is kicked.
	playerMessageBuffer = 16

	// playerWriteTimeout determines how long to
	// wait before removing the player's connection.
	playerWriteTimeout = time.Second * 3
)

// Player represents a player connected to the game server.
// Messages are sent on the messages channel and if the client
// cannot keep up with the messages, closeSlow is called.
type Player struct {
	// uid is the unique id of the player.
	uid string

	// username is the name displayed to others.
	username string

	// team represents the team the player joined.
	team *string

	// messageCh provides a channel the game server use to
	// send messages to the player.
	messageCh chan []byte

	// actionCh provides a channel the player can use to
	// send message to the game server.
	actionCh chan *Action

	// conn is the underlying websocket connection between
	// the player and the game server.
	conn *websocket.Conn

	// closed represents whether or not the websocket
	// connection has been closed.
	closed bool

	// mu ensures closed is thread safe.
	mu sync.Mutex
}

func NewPlayer(uid, username string, conn *websocket.Conn, actionCh chan *Action) *Player {
	return &Player{
		uid:       uid,
		username:  username,
		team:      nil,
		messageCh: make(chan []byte, playerMessageBuffer),
		actionCh:  actionCh,
		conn:      conn,
	}
}

func (p *Player) ReadPump(ctx context.Context) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return net.ErrClosed
	}
	p.mu.Unlock()

	for {
		var action qg.Action
		if err := wsjson.Read(ctx, p.conn, &action); err != nil {
			return err
		}
		p.actionCh <- &Action{
			Action: &action,
			Player: p,
		}
	}
}

func (p *Player) WritePump(ctx context.Context) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return net.ErrClosed
	}
	p.mu.Unlock()

	for {
		select {
		case msg := <-p.messageCh:
			if err := p.write(ctx, msg); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *Player) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	if p.conn != nil {
		p.conn.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
	}
}

func (p *Player) write(ctx context.Context, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, playerWriteTimeout)
	defer cancel()
	return p.conn.Write(ctx, websocket.MessageText, msg)
}
