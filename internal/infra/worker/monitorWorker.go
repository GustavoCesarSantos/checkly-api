package worker

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/nativeSQL"
	urls "GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/handlers"
	unitOfWork_withDB "GustavoCesarSantos/checkly-api/internal/modules/urls/utils/unitOfWork/repositoryUnitOfWork/withDB"
)

type MonitorWorker struct {
	interval    time.Duration
	concurrency int
	monitor     *urls.MonitorUrls
}

func NewMonitorWorker(sqlDB *sql.DB, concurrency int) *MonitorWorker {
	urlRepo := db.NewUrlRepository(sqlDB)
	repoUoW := unitOfWork_withDB.NewRepositoryFactory(sqlDB)
	return &MonitorWorker{
		interval:    1 * time.Minute,
		concurrency: concurrency,
		monitor: urls.NewMonitorUrls(
			application.NewCheckUrl(),
			application.NewEvaluateUrl(),
			application.NewFetchUrls(urlRepo),
			application.NewScheduleNextCheck(),
			application.NewUpdateUrl(urlRepo),
			application.NewUpdateUrlWithOutbox(repoUoW),
		),
	}
}

func (w *MonitorWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	slog.Info("Started", "interval", w.interval.String(), "concurrency", w.concurrency)
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopped")
			return
		case <-ticker.C:
			slog.Info("Tick")
			w.safeProcess(ctx)
		}
	}
}

func (w *MonitorWorker) safeProcess(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from panic", "panic", r)
		}
	}()
	err := w.monitor.Handle(ctx, w.concurrency)
	if err != nil {
		slog.Error("Process error", "error", err)
	} else {
		slog.Info("Cycle completed successfully")
	}
}
