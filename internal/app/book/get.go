package book

import (
	"errors"
	"github.com/BiblioBreeze/core-service/internal/jsonutil"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

type getBookResponse struct {
	ID              uint64  `json:"id"`
	BelongsToUserID uint64  `json:"belongs_to_user_id"`
	Name            string  `json:"name"`
	Author          string  `json:"author"`
	Genre           string  `json:"genre"`
	Description     string  `json:"description"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}

func (s *service) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(1, "Failed to parse id"))
		return
	}

	book, err := s.store.GetBookByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(5, "Book not found"))
			return
		}

		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to get book"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(getBookResponse{
		ID:              book.ID,
		BelongsToUserID: book.BelongsToUserID,
		Name:            book.Name,
		Author:          book.Author,
		Genre:           book.Genre,
		Description:     book.Description,
		Latitude:        book.Latitude,
		Longitude:       book.Longitude,
	}))
}
