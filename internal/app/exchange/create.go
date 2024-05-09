package exchange

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/BiblioBreeze/core-service/internal/app/token"
	"github.com/BiblioBreeze/core-service/internal/jsonutil"
)

type createRequest struct {
	BookID    uint64 `json:"book_id" validate:"required"`
	Condition string `json:"condition" validate:"required"`
}

func (s *service) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createRequest

	statusCode, err := jsonutil.Unmarshal(w, r, &req)
	if err != nil {
		jsonutil.MarshalResponse(w, statusCode, jsonutil.NewError(1, err.Error()))
		return
	}

	book, err := s.store.GetBookByID(ctx, req.BookID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(5, "Book not found"))
			return
		}

		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to get book"))
		return
	}

	userID := token.UserIDFromContext(ctx)

	if book.BelongsToUserID == userID {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(8, "You cannot create exchangeRequest on own book"))
		return
	}

	err = s.store.CreateExchangeRequest(ctx, schema.ExchangeRequest{
		FromUserID: userID,
		BookID:     req.BookID,
		Condition:  req.Condition,
	})
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to create exchangeRequest"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(1))
}
