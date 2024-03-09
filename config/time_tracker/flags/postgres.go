package flags

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostgresFlags struct {
	ConnectionDSN      string        `toml:"dsn"`
	MaxOpenConnections int           `toml:"max-open-connections"`
	ConnectionLifetime time.Duration `toml:"conn-lifetime"`
}

func (f PostgresFlags) Init(_ context.Context) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", f.ConnectionDSN)
	if err != nil {
		return nil, fmt.Errorf("sqlx connect: %v", err)
	}

	db.SetMaxIdleConns(f.MaxOpenConnections)
	db.SetConnMaxLifetime(f.ConnectionLifetime)

	return db, nil
}
