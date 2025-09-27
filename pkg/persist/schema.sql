-- create the quibbble schema
CREATE SCHEMA IF NOT EXISTS quibbble;

-- create the active games table
CREATE TABLE IF NOT EXISTS quibbble.active_games (
    repository VARCHAR NOT NULL;
    tag VARCHAR NOT NULL;
    name VARCHAR NOT NULL;
    snapshot VARBINARY NOT NULL;
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (repository, tag, name)
);

-- create the completed games table
CREATE TABLE IF NOT EXISTS quibbble.completed_games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    repository VARCHAR NOT NULL;
    tag VARCHAR NOT NULL;
    name VARCHAR NOT NULL;
    snapshot VARBINARY NOT NULL;
    updated_at TIMESTAMP NOT NULL
);
