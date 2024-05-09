package user

import (
	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/BiblioBreeze/core-service/internal/app/token"
	"github.com/BiblioBreeze/core-service/internal/jsonutil"
	"net/http"
)

type updateRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

func (s *service) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req updateRequest

	statusCode, err := jsonutil.Unmarshal(w, r, &req)
	if err != nil {
		jsonutil.MarshalResponse(w, statusCode, jsonutil.NewError(1, err.Error()))
		return
	}

	err = s.store.UpdateUser(ctx, schema.User{
		ID:        token.UserIDFromContext(ctx),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to update user"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(1))
}
