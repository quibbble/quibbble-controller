package cockroachdb

import "fmt"

type Config struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

func (c *Config) GetURL() string {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s", c.Username, c.Password, c.Host, c.Database)
	if c.SSLMode != "" {
		url += "?sslmode=" + c.SSLMode
	}
	return url
}
