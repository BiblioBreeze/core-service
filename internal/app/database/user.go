package database

import (
	"context"
	"fmt"
	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateUser(ctx context.Context, user schema.User) error {
	_, err := c.pool.Exec(ctx, `
		INSERT INTO 
		    users(email, first_name, last_name, password)
		VALUES ($1, $2, $3, $4)
	`, user.Email, user.FirstName, user.LastName, user.Password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (schema.User, error) {
	rows, err := c.pool.Query(ctx, `
		SELECT 
			id, email, first_name, last_name, password
		FROM
		    users
		WHERE email = $1
	`, email)
	if err != nil {
		return schema.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[schema.User])
	if err != nil {
		return schema.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
