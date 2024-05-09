package book

import (
	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/BiblioBreeze/core-service/internal/app/token"
	"github.com/BiblioBreeze/core-service/internal/jsonutil"
	"net/http"
)

type createRequest struct {
	Name        string  `json:"name" validate:"required"`
	Author      string  `json:"author" validate:"required"`
	Genre       string  `json:"genre" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Latitude    float64 `json:"latitude" validate:"required,latitude"`
	Longitude   float64 `json:"longitude" validate:"required,longitude"`
}

func (s *service) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createRequest

	statusCode, err := jsonutil.Unmarshal(w, r, &req)
	if err != nil {
		jsonutil.MarshalResponse(w, statusCode, jsonutil.NewError(1, err.Error()))
		return
	}

	err = s.store.CreateBook(ctx, schema.Book{
		BelongsToUserID: token.UserIDFromContext(ctx),
		Name:            req.Name,
		Author:          req.Author,
		Genre:           req.Genre,
		Description:     req.Description,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
	})
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to create book"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(1))
}
