package user

import (
	"errors"
	"github.com/BiblioBreeze/core-service/internal/app/token"
	"github.com/BiblioBreeze/core-service/internal/jsonutil"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

type getUserResponse struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *service) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := token.UserIDFromContext(ctx)
	if idStr := r.URL.Query().Get("id"); idStr != "" {
		idInt, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(1, "Failed to parse id"))
			return
		}

		userID = idInt
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(5, "User not found"))
			return
		}

		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to get user"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(getUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}))
}
