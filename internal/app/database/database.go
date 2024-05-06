package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

func NewClient(pool *pgxpool.Pool) *Client {
	return &Client{pool: pool}
}
