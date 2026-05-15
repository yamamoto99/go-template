package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/yamamoto99/go-template/app/infrastructure/config"
	"github.com/yamamoto99/go-template/app/infrastructure/db"
	"github.com/yamamoto99/go-template/app/internal/handler"
	"github.com/yamamoto99/go-template/app/internal/repository"
	"github.com/yamamoto99/go-template/app/internal/router"
	"github.com/yamamoto99/go-template/app/internal/usecase"
)

const (
	shutdownTimeout = 30 * time.Second
	readTimeout     = 10 * time.Second
	writeTimeout    = 30 * time.Second
	idleTimeout     = 60 * time.Second
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := run(); err != nil {
		slog.Error("startup failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	dbConnection, err := db.NewDB(cfg.DB)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.CloseDB(dbConnection); err != nil {
			slog.Error("close db", "err", err)
		}
	}()

	welcomeRepository := repository.NewWelcomeRepository(dbConnection)
	welcomeUsecase := usecase.NewWelcomeUsecase(welcomeRepository)
	welcomeHandler := handler.NewWelcomeHandler(welcomeUsecase)
	e, err := router.NewRouter(welcomeHandler)
	if err != nil {
		return err
	}
	e.Server.ReadTimeout = readTimeout
	e.Server.WriteTimeout = writeTimeout
	e.Server.IdleTimeout = idleTimeout

	serverErrCh := make(chan error, 1)
	go func() {
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrCh:
		if err != nil {
			return fmt.Errorf("server: %w", err)
		}
		return nil
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
