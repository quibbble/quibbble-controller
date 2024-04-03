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

func (c *Client) GetActiveGame(ctx context.Context, key, id string) (*store.Game, error) {
	if c.pool == nil {
		return nil, store.ErrGameStoreNotEnabled
	}

	sql := `
		SELECT snapshot, updated_at FROM quibbble.active_games
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

func (c *Client) GetActiveGames(ctx context.Context, playerId string) ([]*store.Game, error) {
	if c.pool == nil {
		return nil, store.ErrGameStoreNotEnabled
	}

	sql := `
		SELECT game_key, game_id FROM quibbble.player_games
		WHERE player_id=$1
		ORDER BY updated_at ASC
		LIMIT 20;
	`
	rows, err := c.pool.Query(ctx, sql, playerId)
	if err != nil {
		return nil, err
	}

	games := make([]*store.Game, 0)
	defer rows.Close()
	for rows.Next() {
		var (
			gameKey string
			gameId  string
		)
		err = rows.Scan(&gameKey, &gameId)
		if err != nil {
			return nil, err
		}
		games = append(games, &store.Game{
			Key: gameKey,
			ID:  gameId,
		})
	}

	return games, nil
}

func (c *Client) StoreActiveGame(ctx context.Context, game *store.Game) error {
	if c.pool == nil {
		return store.ErrGameStoreNotEnabled
	}

	// link game with players
	players, err := game.Snapshot.Tags.Players()
	if err != nil {
		return err
	}

	sql := `
		UPSERT INTO quibbble.player_games (player_id, game_key, game_id, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	for _, p := range players {
		for _, player := range p {
			_, err = c.pool.Exec(ctx, sql, player, game.Key, game.ID, game.UpdatedAt)
			if err != nil {
				return store.ErrGameStoreInsert
			}
		}
	}

	// store the game
	sql = `
		UPSERT INTO quibbble.active_games (game_key, game_id, snapshot, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err = c.pool.Exec(ctx, sql, game.Key, game.ID, game.Snapshot.String(), game.UpdatedAt)
	if err != nil {
		return store.ErrGameStoreInsert
	}

	return nil
}

func (c *Client) DeleteActiveGame(ctx context.Context, game *store.Game) error {
	if c.pool == nil {
		return store.ErrGameStoreNotEnabled
	}

	// remove player links
	players, err := game.Snapshot.Tags.Players()
	if err != nil {
		return err
	}

	sql := `
		DELETE FROM quibbble.player_games
		WHERE player_id=$1
		AND game_key=$2
		AND game_id=$3
	`

	for _, p := range players {
		for _, player := range p {
			_, err = c.pool.Exec(ctx, sql, player, game.Key, game.ID)
			if err != nil {
				return store.ErrGameStoreDelete
			}
		}
	}

	// delete the game
	sql = `
		DELETE FROM quibbble.active_games
		WHERE game_key=$1
		AND game_id=$2
	`

	_, err = c.pool.Exec(ctx, sql, game.Key, game.ID)
	if err != nil {
		return store.ErrGameStoreDelete
	}

	return nil
}

func (c *Client) StoreCompletedGame(ctx context.Context, game *store.Game) error {
	if c.pool == nil {
		return store.ErrGameStoreNotEnabled
	}

	sql := `
		INSERT INTO quibbble.completed_games (game_key, snapshot, updated_at)
		VALUES ($1, $2, $3)
	`

	_, err := c.pool.Exec(ctx, sql, game.Key, game.Snapshot.String(), game.UpdatedAt)
	if err != nil {
		return store.ErrGameStoreInsert
	}

	return nil
}

func (c *Client) IncrementGameCount(ctx context.Context, gameKey string) error {
	if c.pool == nil {
		return store.ErrGameStoreNotEnabled
	}

	sql := `
		UPDATE quibbble.created_games
		SET count = count + 1
		WHERE game_key=$1
	`

	_, err := c.pool.Exec(ctx, sql, gameKey)
	if err != nil {
		return store.ErrGameStoreDelete
	}
	return nil
}

func (c *Client) GetStats(ctx context.Context) (*store.Stats, error) {
	if c.pool == nil {
		return nil, store.ErrGameStoreNotEnabled
	}

	sql := `
		WITH A AS (
			SELECT game_key, count AS created_games from quibbble.created_games
		) ,B AS (
			SELECT game_key, COUNT(game_id) AS active_games FROM quibbble.active_games
			GROUP BY game_key
		), C AS (
			SELECT game_key, COUNT(id) AS completed_games FROM quibbble.completed_games
			GROUP BY game_key
		)
		SELECT A.game_key, A.created_games, B.active_games, C.completed_games FROM 
			A JOIN B ON A.game_key = B.game_key
				JOIN C ON B.game_key = C.game_key
	`

	rows, err := c.pool.Query(ctx, sql)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, store.ErrGameStoreNotFound
		}
		return nil, store.ErrGameStoreSelect
	}

	stats := &store.Stats{
		CreatedGameCount:   make(map[string]int),
		ActiveGameCount:    make(map[string]int),
		CompletedGameCount: make(map[string]int),
	}

	var (
		gameKey                                   string
		createdGames, activeGames, completedGames int
	)

	for rows.Next() {
		if err := rows.Scan(&gameKey, &createdGames, &activeGames, &completedGames); err != nil {
			if err == pgx.ErrNoRows {
				return nil, store.ErrGameStoreNotFound
			}
			return nil, store.ErrGameStoreSelect
		}
		stats.CreatedGameCount[gameKey] = createdGames
		stats.ActiveGameCount[gameKey] = activeGames
		stats.CompletedGameCount[gameKey] = completedGames
	}

	return stats, nil
}

func (c *Client) Close(ctx context.Context) error {
	if c.pool == nil {
		return nil
	}
	c.pool.Close()
	return nil
}
