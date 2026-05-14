package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/yamamoto99/go-template/app/infrastructure/db"
	"github.com/yamamoto99/go-template/app/internal/handler"
	"github.com/yamamoto99/go-template/app/internal/repository"
	"github.com/yamamoto99/go-template/app/internal/router"
	"github.com/yamamoto99/go-template/app/internal/usecase"
)

const shutdownTimeout = 30 * time.Second

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("load .env: %w", err)
	}

	dbConnection, err := db.NewDB()
	if err != nil {
		return err
	}
	defer func() {
		if err := db.CloseDB(dbConnection); err != nil {
			log.Println(err)
		}
	}()

	welcomeRepository := repository.NewWelcomeRepository(dbConnection)
	welcomeUsecase := usecase.NewWelcomeUsecase(welcomeRepository)
	welcomeHandler := handler.NewWelcomeHandler(welcomeUsecase)
	e := router.NewRouter(welcomeHandler)

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
		log.Println("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
