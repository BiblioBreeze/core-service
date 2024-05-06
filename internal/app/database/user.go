package database

import (
	"context"
	"fmt"
	"github.com/BiblioBreeze/core-service/internal/app/schema"
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
