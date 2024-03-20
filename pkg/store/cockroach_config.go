package game_store

import "fmt"

type CockroachConfig struct {
	Enabled  bool
	Host     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func (c *CockroachConfig) GetURL() string {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s", c.Username, c.Password, c.Host, c.Database)
	if c.SSLMode != "" {
		url += "?sslmode=" + c.SSLMode
	}
	return url
}
