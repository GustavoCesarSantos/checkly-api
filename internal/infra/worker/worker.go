package worker

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/nativeSQL"
	urls "GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/handlers"
)

type Worker struct {
	interval    time.Duration
	concurrency int
	monitor     *urls.MonitorUrls
}

func NewWorker(sqlDB *sql.DB, concurrency int) *Worker {
	repo := db.NewUrlRepository(sqlDB)
	return &Worker{
		interval:    1 * time.Minute,
		concurrency: concurrency,
		monitor: urls.NewMonitorUrls(
			application.NewCheckUrl(),
			application.NewEvaluateUrl(),
			application.NewScheduleNextCheck(),
			application.NewUpdateUrl(repo),
			repo,
		),
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	slog.Info("worker started", "interval", w.interval.String(), "concurrency", w.concurrency)
	for {
		select {
		case <-ctx.Done():
			slog.Info("worker stopped")
			return
		case <-ticker.C:
			slog.Info("worker tick")
			w.safeProcess(ctx)
		}
	}
}

func (w *Worker) safeProcess(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("worker recovered from panic", "panic", r)
		}
	}()
	err := w.monitor.Handle(ctx, w.concurrency)
	if err != nil {
		slog.Error("worker process error", "error", err)
	} else {
		slog.Info("worker cycle completed successfully")
	}
}
