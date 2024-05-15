-- create the quibbble schema
CREATE SCHEMA IF NOT EXISTS quibbble;

-- create the active games table
CREATE TABLE IF NOT EXISTS quibbble.active_games (
    game_key STRING NOT NULL,
    game_id STRING NOT NULL,
    snapshot STRING NOT NULL,
	updated_at TIMESTAMP,
    CONSTRAINT id PRIMARY KEY (game_key, game_id)
);

-- create the completed games table
CREATE TABLE IF NOT EXISTS quibbble.completed_games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    game_key STRING NOT NULL,
    snapshot STRING NOT NULL, 
    updated_at TIMESTAMP
);

-- create created games table
CREATE TABLE IF NOT EXISTS quibbble.created_games (
    game_key STRING PRIMARY KEY,
    count INT NOT NULL DEFAULT 0
);

-- get a game from the active table
SELECT * FROM quibbble.active_games
WHERE key=$1
AND id=$2

-- insert a game into the active table
INSERT INTO quibbble.active_games (game_key, game_id, snapshot, updated_at)
VALUES ($1, $2, $3, $4)

-- insert a game into the completed table
INSERT INTO quibbble.completed_games (game_key, snapshot, updated_at)
VALUES ($1, $2, $3)

-- increment game count
UPDATE quibbble.created_games
SET count = count + 1
WHERE game_key=$1
