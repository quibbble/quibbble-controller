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

-- create the complete games table
CREATE TABLE IF NOT EXISTS quibbble.complete (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    game_key STRING NOT NULL,
    snapshot STRING NOT NULL, 
    updated_at TIMESTAMP,
);

-- get a game from the active table
SELECT * FROM quibbble.active
WHERE key=$1
AND id=$2

-- insert a game into the active table
INSERT INTO quibbble.games (key, id, snapshot, updated_at)
VALUES ($1, $2, $3, $4)

-- insert a game into the complete table
INSERT INTO quibbble.games (key, snapshot, updated_at)
VALUES ($1, $2, $3)
