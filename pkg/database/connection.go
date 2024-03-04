package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is a Postgres connection pool.
type DB struct {
	Pool *pgxpool.Pool
}

// New returns new Postgres connection pool.
func New(ctx context.Context, url string) (*DB, error) {
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("unable to connection to database: %s", err)
	}

	return &DB{db}, db.Ping(ctx)
}

func (db *DB) Close() {
	db.Pool.Close()
}
