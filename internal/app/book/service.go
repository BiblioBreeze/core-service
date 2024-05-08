package book

import (
	"context"
	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/go-chi/chi/v5"
)

type bookStore interface {
	CreateBook(ctx context.Context, book schema.Book) error
	GetBookByID(ctx context.Context, id uint64) (schema.Book, error)
	ListBooks(ctx context.Context) ([]schema.Book, error)
}

type service struct {
	bookStore bookStore
}

func Mount(r chi.Router, bookStore bookStore) {
	svc := &service{
		bookStore: bookStore,
	}

	svc.routes(r)
}
