package user

import (
	"context"
	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/go-chi/chi/v5"
)

type store interface {
	CreateUser(ctx context.Context, user schema.User) error
	GetUserByEmail(ctx context.Context, email string) (schema.User, error)
}

type tokenService interface {
	Generate(userID uint64) (string, error)
}

type service struct {
	store        store
	tokenService tokenService
}

func Mount(r chi.Router, store store, tokenService tokenService) {
	svc := &service{
		store:        store,
		tokenService: tokenService,
	}

	svc.routes(r)
}
