package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

// Server is a http server.
type Server struct {
	server *http.Server
}

// New returns a new http server.
func New(addr string, router http.Handler) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &Server{
		server: srv,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		slog.Info("stopping server")

		// shutdown our http server, but use new context since the one we
		// waited is cancelled already.
		sCtx, sCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer sCancel()

		if err := s.server.Shutdown(sCtx); err != nil {
			slog.Error("failed to shutdown server", slog.String("error", err.Error()))
		}
	}()

	slog.Info("starting server", slog.String("addr", s.server.Addr))

	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}
