package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

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

func (c *Client) GetUserByID(ctx context.Context, id uint64) (schema.User, error) {
	return c.getUser(ctx, id, "")
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (schema.User, error) {
	return c.getUser(ctx, 0, email)
}

func (c *Client) UpdateUser(ctx context.Context, user schema.User) error {
	_, err := c.pool.Exec(ctx, `
		UPDATE 
		    users 
		SET email = $2, first_name = $3, last_name = $4 
		WHERE id = $1
	`, user.ID, user.Email, user.FirstName, user.LastName)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (c *Client) getUser(ctx context.Context, id uint64, email string) (schema.User, error) {
	whereStatement := "id = $1"
	var whereArg interface{} = id
	if email != "" {
		whereStatement = "email = $1"
		whereArg = email
	}

	rows, err := c.pool.Query(ctx, fmt.Sprintf(`
		SELECT 
			id, email, first_name, last_name, password
		FROM
		    users
		WHERE %s
	`, whereStatement), whereArg)
	if err != nil {
		return schema.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[schema.User])
	if err != nil {
		return schema.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
