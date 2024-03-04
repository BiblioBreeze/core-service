package probes

import (
	"context"
	"core-service/pkg/database"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func SetupFunc(db *database.DB) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/live", liveness())
		r.Get("/ready", readiness(db))
	}
}

func liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func readiness(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		if err := db.Pool.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
