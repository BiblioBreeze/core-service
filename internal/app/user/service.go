package user

import (
	"context"
	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type store interface {
	CreateUser(ctx context.Context, user schema.User) error
	GetUserByID(ctx context.Context, id uint64) (schema.User, error)
	GetUserByEmail(ctx context.Context, email string) (schema.User, error)
	UpdateUser(ctx context.Context, user schema.User) error
}

type tokenService interface {
	Generate(userID uint64) (string, error)
	AuthMiddleware() func(next http.Handler) http.Handler
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
