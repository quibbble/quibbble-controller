-- create the quibbble schema
CREATE SCHEMA IF NOT EXISTS quibbble;

-- create the active games table
CREATE TABLE IF NOT EXISTS quibbble.active (
    game_key STRING NOT NULL,
    game_id STRING NOT NULL,
    snapshot STRING NOT NULL,
	updated_at TIMESTAMP,
    CONSTRAINT id PRIMARY KEY (game_key, game_id)
);

-- create the completed games table
CREATE TABLE IF NOT EXISTS quibbble.completed (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    game_key STRING NOT NULL,
    snapshot STRING NOT NULL, 
    updated_at TIMESTAMP,
);

-- create the player games table
CREATE TABLE IF NOT EXISTS quibbble.player_games (
    player_id STRING NOT NULL,
    game_key STRING NOT NULL,
    game_id STRING NOT NULL,
    updated_at TIMESTAMP,
    CONSTRAINT id PRIMARY KEY (player_id, game_key, game_id)
);

-- create game count table
CREATE TABLE IF NOT EXISTS quibbble.game_count (
    game_key STRING PRIMARY KEY,
    count INT NOT NULL DEFAULT 0,
);

-- get a game from the active table
SELECT * FROM quibbble.active
WHERE key=$1
AND id=$2

-- insert a game into the active table
INSERT INTO quibbble.active (game_key, game_id, snapshot, updated_at)
VALUES ($1, $2, $3, $4)

-- insert a game into the completed table
INSERT INTO quibbble.completed (game_key, snapshot, updated_at)
VALUES ($1, $2, $3)

-- insert a game into the player games table
INSERT INTO quibbble.player_games (player_id, game_key, game_id)
VALUES ($1, $2, $3)

-- increment game count
UPDATE quibbble.game_count
SET count = count + 1
WHERE game_key=$1
