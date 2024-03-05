package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	core_service "github.com/BiblioBreeze/core-service"
	"github.com/BiblioBreeze/core-service/internal/config"
	"github.com/BiblioBreeze/core-service/internal/probes"
	"github.com/BiblioBreeze/core-service/internal/router"
	"github.com/BiblioBreeze/core-service/pkg/database"
	"github.com/BiblioBreeze/core-service/pkg/server"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

var (
	migrate   = flag.Bool("migrate", false, "use app as migrator")
	devLogger = flag.Bool("dev-logger", false, "is dev logger enabled")
)

func main() {
	flag.Parse()

	setUpLogger(*devLogger)

	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to init config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if *migrate {
		if err = runMigrate(cfg.DbDSN); err != nil {
			slog.Error("failed to migrate", slog.String("error", err.Error()))
			os.Exit(1)
		}

		return
	}

	if err = runApp(cfg); err != nil {
		slog.Error("failed to run app", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func runApp(cfg *config.Config) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	db, err := database.New(context.Background(), cfg.DbDSN)
	if err != nil {
		return err
	}

	r := router.New(*devLogger, probes.SetupFunc(db))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// booksService := books.New(...)

	r.Route("/api", func(r chi.Router) {
		// booksService.Routes(r)
	})

	srv := server.New(addr, r)

	slog.Info("server was started successfully", slog.String("addr", addr))

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	if err = srv.Shutdown(ctx); err != nil {
		return errors.New("server shutdown failed")
	}

	db.Close()

	return nil
}

const (
	gooseSQLDriver     = "pgx"
	gooseDialect       = "postgres"
	gooseMigrationsDir = "migrations"
)

func runMigrate(dbDSN string) error {
	var db *sql.DB

	db, err := sql.Open(gooseSQLDriver, dbDSN)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	goose.SetBaseFS(core_service.GetMigrationsFS())

	if err = goose.SetDialect(gooseDialect); err != nil {
		return err
	}

	if err = goose.Up(db, gooseMigrationsDir); err != nil {
		return err
	}

	return nil
}

func setUpLogger(isDevLoggerEnabled bool) {
	var handler slog.Handler

	handler = slog.NewJSONHandler(os.Stdout, nil)
	if isDevLoggerEnabled {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	slog.SetDefault(slog.New(handler))
}
