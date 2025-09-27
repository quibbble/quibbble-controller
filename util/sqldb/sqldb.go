package sqldb

import (
	"context"
	"fmt"
	"time"

	"github.com/quibbble/quibbble-controller/config"
	"github.com/quibbble/quibbble-controller/util/k8s"
	"github.com/upper/db/v4"
	mysqladp "github.com/upper/db/v4/adapter/mysql"
	postgresqladp "github.com/upper/db/v4/adapter/postgresql"
	"k8s.io/client-go/kubernetes"
)

// CreateDBSession creates the dB session
func CreateDBSession(ctx context.Context, clientset kubernetes.Interface, namespace string, persist *config.PersistConfig) (db.Session, error) {
	if persist.PostgreSQL != nil {
		return createPostGresDBSession(ctx, clientset, namespace, persist.PostgreSQL, persist.ConnectionPool)
	} else if persist.MySQL != nil {
		return createMySQLDBSession(ctx, clientset, namespace, persist.MySQL, persist.ConnectionPool)
	}
	return nil, fmt.Errorf("no databases are configured")
}

// CreateDBSessionWithCreds creates a database session using direct username and password
func CreateDBSessionWithCreds(persist *config.PersistConfig, username, password string) (db.Session, error) {
	if persist.PostgreSQL != nil {
		return createPostGresDBSessionWithCreds(persist.PostgreSQL, persist.ConnectionPool, username, password)
	} else if persist.MySQL != nil {
		return createMySQLDBSessionWithCreds(persist.MySQL, persist.ConnectionPool, username, password)
	}
	return nil, fmt.Errorf("no databases are configured")
}

// createPostGresDBSession creates postgresDB session
func createPostGresDBSession(ctx context.Context, clientset kubernetes.Interface, namespace string, cfg *config.PostgreSQLConfig, persistPool *config.ConnectionPool) (db.Session, error) {
	userNameByte, err := k8s.GetSecrets(ctx, clientset, namespace, cfg.UsernameSecret.Name, cfg.UsernameSecret.Key)
	if err != nil {
		return nil, err
	}
	passwordByte, err := k8s.GetSecrets(ctx, clientset, namespace, cfg.PasswordSecret.Name, cfg.PasswordSecret.Key)
	if err != nil {
		return nil, err
	}

	return createPostGresDBSessionWithCreds(cfg, persistPool, string(userNameByte), string(passwordByte))
}

// createMySQLDBSession creates Mysql DB session
func createMySQLDBSession(ctx context.Context, clientset kubernetes.Interface, namespace string, cfg *config.MySQLConfig, persistPool *config.ConnectionPool) (db.Session, error) {
	userNameByte, err := k8s.GetSecrets(ctx, clientset, namespace, cfg.UsernameSecret.Name, cfg.UsernameSecret.Key)
	if err != nil {
		return nil, err
	}
	passwordByte, err := k8s.GetSecrets(ctx, clientset, namespace, cfg.PasswordSecret.Name, cfg.PasswordSecret.Key)
	if err != nil {
		return nil, err
	}

	return createMySQLDBSessionWithCreds(cfg, persistPool, string(userNameByte), string(passwordByte))
}

// createPostGresDBSessionWithCreds creates postgresDB session with direct credentials
func createPostGresDBSessionWithCreds(cfg *config.PostgreSQLConfig, persistPool *config.ConnectionPool, username, password string) (db.Session, error) {
	settings := postgresqladp.ConnectionURL{
		User:     username,
		Password: password,
		Host:     cfg.GetHostname(),
		Database: cfg.Database,
	}

	if cfg.SSL {
		if cfg.SSLMode != "" {
			options := map[string]string{
				"sslmode": cfg.SSLMode,
			}
			settings.Options = options
		}
	}

	session, err := postgresqladp.Open(settings)
	if err != nil {
		return nil, err
	}
	session = ConfigureDBSession(session, persistPool)
	return session, nil
}

// createMySQLDBSessionWithCreds creates MySQL DB session with direct credentials
func createMySQLDBSessionWithCreds(cfg *config.MySQLConfig, persistPool *config.ConnectionPool, username, password string) (db.Session, error) {
	session, err := mysqladp.Open(mysqladp.ConnectionURL{
		User:     username,
		Password: password,
		Host:     cfg.GetHostname(),
		Database: cfg.Database,
		Options:  cfg.Options,
	})
	if err != nil {
		return nil, err
	}
	session = ConfigureDBSession(session, persistPool)
	// this is needed to make MySQL run in a Golang-compatible UTF-8 character set.
	_, err = session.SQL().Exec("SET NAMES 'utf8mb4'")
	if err != nil {
		return nil, err
	}
	_, err = session.SQL().Exec("SET CHARACTER SET utf8mb4")
	if err != nil {
		return nil, err
	}
	return session, nil
}

// ConfigureDBSession configures the DB session
func ConfigureDBSession(session db.Session, dbPool *config.ConnectionPool) db.Session {
	if dbPool != nil {
		session.SetMaxOpenConns(dbPool.MaxOpenConns)
		session.SetMaxIdleConns(dbPool.MaxIdleConns)
		session.SetConnMaxLifetime(time.Duration(dbPool.ConnMaxLifetime))
	}
	return session
}
