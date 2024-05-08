package token

import (
	"context"
	"github.com/BiblioBreeze/core-service/internal/jsonutil"
	"net/http"
	"strings"
)

const (
	authorizationHeaderName = "Authorization"
)

type authContextKey struct{}

var (
	unauthorizedJsonErr = jsonutil.NewError(7, "Not authorized")
)

func (s *Service) AuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(authorizationHeaderName)
			t := strings.Split(authHeader, " ")

			if len(t) == 2 {
				userID, err := s.Validate(t[1])
				if err != nil || userID == 0 {
					jsonutil.MarshalResponse(w, http.StatusUnauthorized, unauthorizedJsonErr)
					return
				}

				next.ServeHTTP(
					w,
					r.WithContext(context.WithValue(r.Context(), authContextKey{}, userID)),
				)

				return
			}

			jsonutil.MarshalResponse(w, http.StatusUnauthorized, unauthorizedJsonErr)
			return
		})
	}
}

func UserIDFromContext(ctx context.Context) uint64 {
	val, _ := ctx.Value(authContextKey{}).(uint64)
	return val
}
