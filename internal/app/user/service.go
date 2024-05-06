package user

import (
	"context"
	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/go-chi/chi/v5"
)

type userStore interface {
	CreateUser(ctx context.Context, user schema.User) error
}

type service struct {
	userStore userStore
}

func Mount(r chi.Router, userStore userStore) {
	svc := &service{
		userStore: userStore,
	}

	svc.routes(r)
}
