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

	"GustavoCesarSantos/checkly-api/internal/shared/configs"
	"GustavoCesarSantos/checkly-api/internal/shared/logger"
)

func Server(db *sql.DB, wg *sync.WaitGroup) error {
	serverConfigs := configs.LoadServerConfig()
	port := serverConfigs.Port
	srvLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	routes := routes(db)
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        routes,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       slog.NewLogLogger(srvLogger.Handler(), slog.LevelError),
	}
	shutdownError := make(chan error)
	go handleShutdown(srv, shutdownError, wg)
	logger.Info(
		fmt.Sprintf("starting server on :%d", port),
		"server.go",
		"Server",
	)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server: %w", err)
	}
	err = <-shutdownError
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}
	logger.Info(
		"server stopped",
		"server.go",
		"Server",
		"addr", srv.Addr,
	)
	return nil
}

func handleShutdown(srv *http.Server, shutdownError chan error, wg *sync.WaitGroup) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	logger.Info(
		"shutting down server",
		"server.go",
		"handleShutdown",
		"signal", s.String(),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		shutdownError <- err
	}
	logger.Info(
		"completing background tasks",
		"server.go",
		"handleShutdown",
		"addr", srv.Addr,
	)
	wg.Wait()
	shutdownError <- nil
}
