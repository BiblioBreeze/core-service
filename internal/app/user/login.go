package user

import (
	"errors"
	"github.com/BiblioBreeze/core-service/internal/jsonutil"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *service) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req loginRequest

	statusCode, err := jsonutil.Unmarshal(w, r, &req)
	if err != nil {
		jsonutil.MarshalResponse(w, statusCode, jsonutil.NewError(1, err.Error()))
		return
	}

	user, err := s.userStore.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(5, "User not exists or password incorrect"))
			return
		}

		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to get user"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(5, "User not exists or password incorrect"))
		return
	}

	token, err := s.tokenService.Generate(user.ID)
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to generate token"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(loginResponse{
		AccessToken: token,
	}))
}
