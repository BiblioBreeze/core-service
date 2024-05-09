package user

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/BiblioBreeze/core-service/internal/jsonutil"
)

type registerRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (s *service) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerRequest

	statusCode, err := jsonutil.Unmarshal(w, r, &req)
	if err != nil {
		jsonutil.MarshalResponse(w, statusCode, jsonutil.NewError(1, err.Error()))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(3, "Failed to hash password"))
		return
	}

	err = s.store.CreateUser(
		ctx,
		schema.User{
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Password:  string(hashedPassword),
		},
	)
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to create user"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(1))
}
