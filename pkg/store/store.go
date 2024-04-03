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
	Snapshot  *qgn.Snapshot `json:"snapshot,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitempty"`
}

type Stats struct {
	CreatedGameCount   map[string]int `json:"created_game_count"`
	ActiveGameCount    map[string]int `json:"active_game_count"`
	CompletedGameCount map[string]int `json:"completed_game_count"`
}

// GameStore stores games into long term storage
type GameStore interface {
	// GetActiveGame retrieves game data for a game
	GetActiveGame(ctx context.Context, key, id string) (*Game, error)

	// GetActiveGames retrieves a list of key-ids tied to a given player
	GetActiveGames(ctx context.Context, playerId string) ([]*Game, error)

	// StoreActiveGame stores a game in active storage
	StoreActiveGame(ctx context.Context, game *Game) error

	// DeleteActive removes a game from active storage
	DeleteActiveGame(ctx context.Context, game *Game) error

	// StoreCompleted stores a game in completed storage
	StoreCompletedGame(ctx context.Context, game *Game) error

	// IncrementGameCount increments the number of games created by 1
	IncrementGameCount(ctx context.Context, gameKey string) error

	// GetStats retrieves created, active, and completed game counts
	GetStats(ctx context.Context) (*Stats, error)

	// Close closes the store
	Close(ctx context.Context) error
}
