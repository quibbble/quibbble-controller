package cockroachdb

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	"github.com/quibbble/quibbble-controller/pkg/store"
)

type Client struct {
	pool *pgxpool.Pool
}

func NewClient(config *Config) (*Client, error) {
	if !config.Enabled {
		return &Client{}, nil
	}

	cfg, err := pgxpool.ParseConfig(config.GetURL())
	if err != nil {
		return nil, store.ErrGameStoreConnection
	}

	cfg.MaxConns = 3
	cfg.MinConns = 0
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = time.Minute * 30
	cfg.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, store.ErrGameStoreConnection
	}

	return &Client{
		pool: pool,
	}, nil
}

func (c *Client) GetGame(ctx context.Context, key, id string) (*store.Game, error) {
	if c.pool == nil {
		return nil, store.ErrGameStoreNotEnabled
	}

	sql := `
		SELECT snapshot, updated_at FROM quibbble.active
		WHERE game_key=$1
		AND game_id=$2
	`
	row := c.pool.QueryRow(ctx, sql, key, id)

	var (
		raw       string
		updatedAt time.Time
	)

	if err := row.Scan(&raw, &updatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, store.ErrGameStoreNotFound
		}
		return nil, store.ErrGameStoreSelect
	}

	snapshot, err := qgn.Parse(raw)
	if err != nil {
		return nil, err
	}

	return &store.Game{
		Key:       key,
		ID:        id,
		Snapshot:  snapshot,
		UpdatedAt: updatedAt,
	}, nil
}

func (c *Client) GetStats(ctx context.Context) (*store.Stats, error) {
	if c.pool == nil {
		return nil, store.ErrGameStoreNotEnabled
	}

	sql := `
		WITH A AS (
			SELECT game_key, COUNT(game_id) AS active_games FROM quibbble.active
			GROUP BY game_key
		), B AS (
			SELECT game_key, COUNT(id) AS complete_games FROM quibbble.complete
			GROUP BY game_key
		)
		SELECT A.game_key, A.active_games, B.complete_games FROM A JOIN B ON A.game_key = B.game_key
	`

	rows, err := c.pool.Query(ctx, sql)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, store.ErrGameStoreNotFound
		}
		return nil, store.ErrGameStoreSelect
	}

	stats := &store.Stats{
		ActiveGames:   make(map[string]int),
		CompleteGames: make(map[string]int),
	}

	var (
		gameKey                    string
		activeGames, completeGames int
	)

	for rows.Next() {
		if err := rows.Scan(&gameKey, &activeGames, &completeGames); err != nil {
			if err == pgx.ErrNoRows {
				return nil, store.ErrGameStoreNotFound
			}
			return nil, store.ErrGameStoreSelect
		}
		stats.ActiveGames[gameKey] = activeGames
		stats.CompleteGames[gameKey] = completeGames
	}

	return stats, nil
}

func (c *Client) StoreActive(ctx context.Context, game *store.Game) error {
	if c.pool == nil {
		return store.ErrGameStoreNotEnabled
	}

	sql := `
		UPSERT INTO quibbble.active (game_key, game_id, snapshot, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := c.pool.Exec(ctx, sql, game.Key, game.ID, game.Snapshot.String(), game.UpdatedAt)
	if err != nil {
		return store.ErrGameStoreInsert
	}

	return nil
}

func (c *Client) DeleteActive(ctx context.Context, game *store.Game) error {
	if c.pool == nil {
		return store.ErrGameStoreNotEnabled
	}

	sql := `
		DELETE FROM quibbble.active
		WHERE game_key=$1
		AND game_id=$2
	`

	_, err := c.pool.Exec(ctx, sql, game.Key, game.ID)
	if err != nil {
		return store.ErrGameStoreDelete
	}

	return nil
}

func (c *Client) StoreComplete(ctx context.Context, game *store.Game) error {
	if c.pool == nil {
		return store.ErrGameStoreNotEnabled
	}

	sql := `
		INSERT INTO quibbble.complete (game_key, snapshot, updated_at)
		VALUES ($1, $2, $3)
	`

	_, err := c.pool.Exec(ctx, sql, game.Key, game.Snapshot.String(), game.UpdatedAt)
	if err != nil {
		return store.ErrGameStoreInsert
	}

	return nil
}

func (c *Client) Close(ctx context.Context) error {
	if c.pool == nil {
		return nil
	}
	c.pool.Close()
	return nil
}
