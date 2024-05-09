package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/BiblioBreeze/core-service/internal/app/schema"
)

func (c *Client) CreateExchangeRequest(ctx context.Context, exchangeRequest schema.ExchangeRequest) error {
	_, err := c.pool.Exec(ctx, `
		INSERT INTO 
		    exchange_requests(from_user_id, book_id, condition, exchanged)
		VALUES ($1, $2, $3, $4)
	`, exchangeRequest.FromUserID, exchangeRequest.BookID, exchangeRequest.Condition, exchangeRequest.Exchanged)
	if err != nil {
		return fmt.Errorf("failed to create exchangeRequest: %w", err)
	}

	return nil
}

func (c *Client) ListExchangeRequests(ctx context.Context, userID uint64) ([]schema.ExchangeRequest, error) {
	rows, err := c.pool.Query(ctx, `
		SELECT
    		er.id, er.from_user_id, er.book_id, er.condition, er.exchanged
		FROM
    		exchange_requests as er
		JOIN books b on b.id = er.book_id
		WHERE er.from_user_id = $1 OR b.belongs_to_user_id = $1
		ORDER BY er.id DESC;
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query exchangeRequests: %w", err)
	}

	exchangeRequests, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[schema.ExchangeRequest])
	if err != nil {
		return nil, fmt.Errorf("failed to collect exchangeRequests: %w", err)
	}

	return exchangeRequests, nil
}
