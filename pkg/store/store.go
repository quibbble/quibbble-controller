package store

import (
	"context"
	"fmt"
	"time"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

var (
	ErrGameStoreNotEnabled = fmt.Errorf("game store is not enabled")
	ErrGameStoreNotFound   = fmt.Errorf("no game found in game store")
	ErrGameStoreConnection = fmt.Errorf("failed to connect to game store")
	ErrGameStoreSelect     = fmt.Errorf("failed to select from game store")
	ErrGameStoreInsert     = fmt.Errorf("failed to insert into game store")
	ErrGameStoreDelete     = fmt.Errorf("failed to delete from game store")
)

type Game struct {
	Key       string        `json:"key"`
	ID        string        `json:"id"`
	Snapshot  *qgn.Snapshot `json:"snapshot"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type Stats struct {
	CreatedGameCount  map[string]int `json:"created_game_count"`
	CompleteGameCount map[string]int `json:"complete_game_count"`
}

// GameStore stores games into long term storage
type GameStore interface {
	GetGame(ctx context.Context, key, id string) (*Game, error)
	GetStats(ctx context.Context) (*Stats, error)
	StoreActive(ctx context.Context, game *Game) error
	DeleteActive(ctx context.Context, game *Game) error
	StoreComplete(ctx context.Context, game *Game) error
	Close(ctx context.Context) error
}
