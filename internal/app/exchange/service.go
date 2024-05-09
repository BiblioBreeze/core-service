package exchange

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/BiblioBreeze/core-service/internal/app/schema"
)

type store interface {
	GetBookByID(ctx context.Context, id uint64) (schema.Book, error)

	CreateExchangeRequest(ctx context.Context, exchangeRequest schema.ExchangeRequest) error
	ListExchangeRequests(ctx context.Context, userID uint64) ([]schema.ExchangeRequest, error)
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
