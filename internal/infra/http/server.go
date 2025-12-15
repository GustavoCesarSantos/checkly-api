package http

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"GustavoCesarSantos/checkly-api/internal/shared/configs"
)

func Server(db *sql.DB) error {
	loadEnvErr := godotenv.Load()
	if loadEnvErr != nil {
		slog.Error("failed to load .env file", "error", loadEnvErr)
		return loadEnvErr
	}
	serverConfigs := configs.LoadServerConfig()
	routes := routes(db)
	port := serverConfigs.Port
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        routes,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	shutdownError := make(chan error)
	var wg sync.WaitGroup
	go handleShutdown(srv, shutdownError, &wg)
	slog.Info(fmt.Sprintf("starting server on :%d", port))
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	slog.Info("stopped server", "addr", srv.Addr)
	return nil
}

func handleShutdown(srv *http.Server, shutdownError chan error, wg *sync.WaitGroup) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	slog.Info("shutting down server", "signal", s.String())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		shutdownError <- err
	}
	slog.Info("completing background tasks", "addr", srv.Addr)
	wg.Wait()
	shutdownError <- nil
}
