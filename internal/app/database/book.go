package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/BiblioBreeze/core-service/internal/app/schema"
)

func (c *Client) CreateBook(ctx context.Context, book schema.Book) error {
	_, err := c.pool.Exec(ctx, `
		INSERT INTO 
		    books(belongs_to_user_id, name, author, genre, description, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, book.BelongsToUserID, book.Name, book.Author, book.Genre, book.Description, book.Latitude, book.Longitude)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}

func (c *Client) ListBooks(ctx context.Context) ([]schema.Book, error) {
	rows, err := c.pool.Query(ctx, `
		SELECT
    		b.id, b.belongs_to_user_id, b.name, b.author, b.genre, b.description, b.latitude, b.longitude
		FROM books as b
    	LEFT JOIN exchange_requests er on b.id = er.book_id and er.exchanged
		WHERE er.exchanged IS NULL;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query books: %w", err)
	}

	books, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[schema.Book])
	if err != nil {
		return nil, fmt.Errorf("failed to collect books: %w", err)
	}

	return books, nil
}

func (c *Client) GetBookByID(ctx context.Context, id uint64) (schema.Book, error) {
	rows, err := c.pool.Query(ctx, `
		SELECT
    		id, belongs_to_user_id, name, author, genre, description, latitude, longitude
		FROM 
		    books
		WHERE 
		    id = $1
	`, id)
	if err != nil {
		return schema.Book{}, fmt.Errorf("failed to get book: %w", err)
	}

	book, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[schema.Book])
	if err != nil {
		return schema.Book{}, fmt.Errorf("failed to collect book: %w", err)
	}

	return book, nil
}
