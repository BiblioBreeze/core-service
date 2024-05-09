package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	core_service "github.com/BiblioBreeze/core-service"
	bookService "github.com/BiblioBreeze/core-service/internal/app/book"
	"github.com/BiblioBreeze/core-service/internal/app/database"
	exchangeService "github.com/BiblioBreeze/core-service/internal/app/exchange"
	tokenService "github.com/BiblioBreeze/core-service/internal/app/token"
	userService "github.com/BiblioBreeze/core-service/internal/app/user"
	"github.com/BiblioBreeze/core-service/internal/config"
	"github.com/BiblioBreeze/core-service/internal/router"
	"github.com/BiblioBreeze/core-service/pkg/server"
	"github.com/BiblioBreeze/core-service/pkg/signal"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
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
	ctx := signal.Context()

	db, err := connectToPostgres(ctx, cfg.DbDSN)
	if err != nil {
		return err
	}

	defer db.Close()

	r := router.New(*devLogger)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	dbClient := database.New(db)
	tokenSvc := tokenService.New(cfg.JWTSigningKey)

	r.Route("/api", func(r chi.Router) {
		userService.Mount(
			r,
			dbClient,
			tokenSvc,
		)
		bookService.Mount(
			r.With(tokenSvc.AuthMiddleware()),
			dbClient,
		)
		exchangeService.Mount(
			r.With(tokenSvc.AuthMiddleware()),
			dbClient,
		)
	})

	srv := server.New(addr, r)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return srv.Run(gCtx)
	})

	if err = g.Wait(); err != nil {
		slog.Error("an un-recoverable error occurred", slog.String("error", err.Error()))
		return err
	}

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

func connectToPostgres(ctx context.Context, url string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("unable to connection to database: %s", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return db, nil
}
