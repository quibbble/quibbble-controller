package game_store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

type CockroachClient struct {
	pool *pgxpool.Pool
}

func NewCockroachClient(config *CockroachConfig) (*CockroachClient, error) {
	if !config.Enabled {
		return &CockroachClient{}, nil
	}

	cfg, err := pgxpool.ParseConfig(config.GetURL())
	if err != nil {
		return nil, ErrGameStoreConnection
	}

	cfg.MaxConns = 3
	cfg.MinConns = 0
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = time.Minute * 30
	cfg.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, ErrGameStoreConnection
	}

	return &CockroachClient{
		pool: pool,
	}, nil
}

func (c *CockroachClient) GetGame(ctx context.Context, key, id string) (*Game, error) {
	if c.pool == nil {
		return nil, ErrGameStoreNotEnabled
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
			return nil, ErrGameStoreNotFound
		}
		return nil, ErrGameStoreSelect
	}

	snapshot, err := qgn.Parse(raw)
	if err != nil {
		return nil, err
	}

	return &Game{
		Key:       key,
		ID:        id,
		Snapshot:  snapshot,
		UpdatedAt: updatedAt,
	}, nil
}

func (c *CockroachClient) GetStats(ctx context.Context) (*Stats, error) {
	if c.pool == nil {
		return nil, ErrGameStoreNotEnabled
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
			return nil, ErrGameStoreNotFound
		}
		return nil, ErrGameStoreSelect
	}

	stats := &Stats{
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
				return nil, ErrGameStoreNotFound
			}
			return nil, ErrGameStoreSelect
		}
		stats.ActiveGames[gameKey] = activeGames
		stats.CompleteGames[gameKey] = completeGames
	}

	return stats, nil
}

func (c *CockroachClient) StoreActive(ctx context.Context, game *Game) error {
	if c.pool == nil {
		return ErrGameStoreNotEnabled
	}

	sql := `
		UPSERT INTO quibbble.active (game_key, game_id, snapshot, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := c.pool.Exec(ctx, sql, game.Key, game.ID, game.Snapshot.String(), game.UpdatedAt)
	if err != nil {
		return ErrGameStoreInsert
	}

	return nil
}

func (c *CockroachClient) DeleteActive(ctx context.Context, game *Game) error {
	if c.pool == nil {
		return ErrGameStoreNotEnabled
	}

	sql := `
		DELETE FROM quibbble.active
		WHERE game_key=$1
		AND game_id=$2
	`

	_, err := c.pool.Exec(ctx, sql, game.Key, game.ID)
	if err != nil {
		return ErrGameStoreDelete
	}

	return nil
}

func (c *CockroachClient) StoreComplete(ctx context.Context, game *Game) error {
	if c.pool == nil {
		return ErrGameStoreNotEnabled
	}

	sql := `
		INSERT INTO quibbble.complete (game_key, snapshot, updated_at)
		VALUES ($1, $2, $3)
	`

	_, err := c.pool.Exec(ctx, sql, game.Key, game.Snapshot.String(), game.UpdatedAt)
	if err != nil {
		return ErrGameStoreInsert
	}

	return nil
}

func (c *CockroachClient) Close(ctx context.Context) error {
	if c.pool == nil {
		return nil
	}
	c.pool.Close()
	return nil
}
