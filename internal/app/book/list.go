package book

import (
	"net/http"

	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/BiblioBreeze/core-service/internal/jsonutil"
	"github.com/BiblioBreeze/core-service/pkg/utils/sliceutils"
)

type listBook struct {
	ID        uint64  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (s *service) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	books, err := s.store.ListBooks(ctx)
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to get books"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(sliceutils.ConvertFunc(books, bindListBookFromModel)))
}

func bindListBookFromModel(book schema.Book) listBook {
	return listBook{
		ID:        book.ID,
		Name:      book.Name,
		Latitude:  book.Latitude,
		Longitude: book.Longitude,
	}
}
