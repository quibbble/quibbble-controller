package config

import (
	"fmt"
	"time"

	apiv1 "k8s.io/api/core/v1"
)

// PersistConfig contains database configuration settings
type PersistConfig struct {
	// PostgreSQL configuration for PostgreSQL database, don't use MySQL at the same time
	PostgreSQL *PostgreSQLConfig `json:"postgresql,omitempty"`
	// MySQL configuration for MySQL database, don't use PostgreSQL at the same time
	MySQL *MySQLConfig `json:"mysql,omitempty"`
	// Pooled connection settings for all types of database connections
	ConnectionPool *ConnectionPool `json:"connectionPool,omitempty"`
}

// PostgreSQLConfig contains PostgreSQL-specific database configuration
type PostgreSQLConfig struct {
	DatabaseConfig
	// SSL enables SSL connection to the database
	SSL bool `json:"ssl,omitempty"`
	// SSLMode specifies the SSL mode (disable, require, verify-ca, verify-full)
	SSLMode string `json:"sslMode,omitempty"`
}

// MySQLConfig contains MySQL-specific database configuration
type MySQLConfig struct {
	DatabaseConfig
	// Options contains additional MySQL connection options
	Options map[string]string `json:"options,omitempty"`
}

// ConnectionPool contains database connection pool settings
type ConnectionPool struct {
	// MaxIdleConns sets the maximum number of idle connections in the pool
	MaxIdleConns int `json:"maxIdleConns,omitempty"`
	// MaxOpenConns sets the maximum number of open connections to the database
	MaxOpenConns int `json:"maxOpenConns,omitempty"`
	// ConnMaxLifetime sets the maximum amount of time a connection may be reused
	ConnMaxLifetime time.Duration `json:"connMaxLifetime,omitempty"`
}

// DatabaseConfig contains common database connection settings
type DatabaseConfig struct {
	// Host is the database server hostname
	Host string `json:"host"`
	// Port is the database server port
	Port int `json:"port,omitempty"`
	// Database is the name of the database to connect to
	Database string `json:"database"`
	// TableName is the name of the table to use, must be set
	TableName string `json:"tableName,omitempty"`
	// UsernameSecret references a secret containing the database username
	UsernameSecret apiv1.SecretKeySelector `json:"userNameSecret,omitempty"`
	// PasswordSecret references a secret containing the database password
	PasswordSecret apiv1.SecretKeySelector `json:"passwordSecret,omitempty"`
}

func (c DatabaseConfig) GetHostname() string {
	if c.Port == 0 {
		return c.Host
	}
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}

func (p *PersistConfig) GetTableName() (string, error) {
	var table string
	if p.PostgreSQL != nil {
		table = p.PostgreSQL.TableName
	} else if p.MySQL != nil {
		table = p.MySQL.TableName
	}
	if table == "" {
		return "", fmt.Errorf("tableName is missing")
	}
	return table, nil
}
