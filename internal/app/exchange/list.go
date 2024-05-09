package exchange

import (
	"net/http"

	"github.com/BiblioBreeze/core-service/internal/app/schema"
	"github.com/BiblioBreeze/core-service/internal/app/token"
	"github.com/BiblioBreeze/core-service/internal/jsonutil"
	"github.com/BiblioBreeze/core-service/pkg/utils/sliceutils"
)

type listExchangeRequest struct {
	ID         uint64 `json:"id"`
	FromUserID uint64 `json:"from_user_id"`
	BookID     uint64 `json:"book_id"`
	Condition  string `json:"condition"`
	Exchanged  bool   `json:"exchanged"`
}

func (s *service) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchangeRequests, err := s.store.ListExchangeRequests(ctx, token.UserIDFromContext(ctx))
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to get exchangeRequests"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(sliceutils.ConvertFunc(exchangeRequests, bindListExchangeRequestFromModel)))
}

func bindListExchangeRequestFromModel(exchangeRequest schema.ExchangeRequest) listExchangeRequest {
	return listExchangeRequest{
		ID:         exchangeRequest.ID,
		FromUserID: exchangeRequest.FromUserID,
		BookID:     exchangeRequest.BookID,
		Condition:  exchangeRequest.Condition,
		Exchanged:  exchangeRequest.Exchanged,
	}
}
