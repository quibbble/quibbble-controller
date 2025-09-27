package persist

import "time"

const (
	ActiveGamesTable    = "active_games"
	CompletedGamesTable = "completed_games"

	RepositoryColumn = "repository"
	TagColumn        = "tag"
	NameColumn       = "name"
	SnapshotColumn   = "snapshot"
	UpdatedAtColumn  = "updated_at"
)

type Game struct {
	Repository string    `db:"repository"`
	Tag        string    `db:"tag"`
	Name       string    `db:"name"`
	Snapshot   []byte    `db:"snapshot"`
	UpdatedAt  time.Time `db:"updated_at"`
}
