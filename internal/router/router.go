package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	slogchi "github.com/samber/slog-chi"
	"github.com/yarlson/chiprom"

	"github.com/BiblioBreeze/core-service/internal/jsonutil"
)

const (
	serviceName = "core-service"
)

// New returns new chi.Mux.
func New(isDevLoggerEnabled bool) *chi.Mux {
	r := chi.NewRouter()

	if isDevLoggerEnabled {
		r.Use(loggerMiddleware())
	}

	r.Use(chiprom.NewMiddleware(serviceName))
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		jsonutil.MarshalResponse(w, http.StatusNotFound, jsonutil.NewError(3, "API method not found"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		jsonutil.MarshalResponse(w, http.StatusMethodNotAllowed, jsonutil.NewError(3, "HTTP method not allowed"))
	})

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}

func loggerMiddleware() func(http.Handler) http.Handler {
	return slogchi.NewWithConfig(
		slog.Default(),
		slogchi.Config{
			DefaultLevel:     slog.LevelInfo,
			ClientErrorLevel: slog.LevelWarn,
			ServerErrorLevel: slog.LevelError,
			WithRequestID:    true,
			WithSpanID:       true,
			WithTraceID:      true,
		},
	)
}
