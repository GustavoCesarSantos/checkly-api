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

type Worker struct {
	interval    time.Duration
	concurrency int
	monitor     *urls.MonitorUrls
}

func NewWorker(sqlDB *sql.DB, concurrency int) *Worker {
	urlRepo := db.NewUrlRepository(sqlDB)
	repoUoW := unitOfWork_withDB.NewRepositoryFactory(sqlDB)
	return &Worker{
		interval:    1 * time.Minute,
		concurrency: concurrency,
		monitor: urls.NewMonitorUrls(
			application.NewCheckUrl(),
			application.NewEvaluateUrl(),
			application.NewScheduleNextCheck(),
			application.NewUpdateUrl(urlRepo),
			application.NewUpdateUrlWithOutbox(repoUoW),
			urlRepo,
		),
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	slog.Info("[WORKER] Started", "interval", w.interval.String(), "concurrency", w.concurrency)
	for {
		select {
		case <-ctx.Done():
			slog.Info("[WORKER] Stopped")
			return
		case <-ticker.C:
			slog.Info("[WORKER] Tick")
			w.safeProcess(ctx)
		}
	}
}

func (w *Worker) safeProcess(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("[WORKER] Recovered from panic", "panic", r)
		}
	}()
	err := w.monitor.Handle(ctx, w.concurrency)
	if err != nil {
		slog.Error("[WORKER] Process error", "error", err)
	} else {
		slog.Info("[WORKER] Cycle completed successfully")
	}
}
