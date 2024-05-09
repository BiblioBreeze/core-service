package book

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/BiblioBreeze/core-service/internal/app/schema"
)

type store interface {
	CreateBook(ctx context.Context, book schema.Book) error
	GetBookByID(ctx context.Context, id uint64) (schema.Book, error)
	ListBooks(ctx context.Context) ([]schema.Book, error)
}

type service struct {
	store store
}

func Mount(r chi.Router, store store) {
	svc := &service{
		store: store,
	}

	svc.routes(r)
}
