package worker

import (
	"context"
	"database/sql"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/nativeSQL"
	urls "GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/handlers"
	unitOfWork_withDB "GustavoCesarSantos/checkly-api/internal/modules/urls/utils/unitOfWork/repositoryUnitOfWork/withDB"
	"GustavoCesarSantos/checkly-api/internal/shared/logger"
)

type MonitorWorker struct {
	interval    time.Duration
	concurrency int
	monitor     *urls.MonitorUrls
	running     chan struct{}
}

func NewMonitorWorker(sqlDB *sql.DB, concurrency int) *MonitorWorker {
	urlRepo := db.NewUrlRepository(sqlDB)
	repoUoW := unitOfWork_withDB.NewRepositoryFactory(sqlDB)
	return &MonitorWorker{
		interval:    1 * time.Minute,
		concurrency: concurrency,
		running:     make(chan struct{}, 1),
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
	logger.InfoContext(
		ctx,
		"Monitor worker started",
		"monitorWorker.go",
		"Start",
		"interval", w.interval.String(),
		"concurrency", w.concurrency,
	)
	for {
		select {
		case <-ctx.Done():
			w.waitForRunning()
			logger.InfoContext(
				ctx,
				"Monitor worker stopped",
				"monitorWorker.go",
				"Start",
			)
			return
		case <-ticker.C:
			select {
			case w.running <- struct{}{}:
				logger.InfoContext(
					ctx,
					"Monitor worker tick",
					"monitorWorker.go",
					"Start",
				)
				go w.runCycle(ctx)
			default:
				logger.WarnContext(
					ctx,
					"Previous monitor worker cycle still running, skipping this tick",
					"monitorWorker.go",
					"Start",
				)
			}
		}
	}
}

func (w *MonitorWorker) waitForRunning() {
	select {
	case w.running <- struct{}{}:
		<-w.running
	default:
	}
}

func (w *MonitorWorker) runCycle(ctx context.Context) {
	defer func() {
		<-w.running
		if r := recover(); r != nil {
			logger.Warn(
				"recovered from panic",
				"monitorWorker.go",
				"runCycle",
				"panic", r,
			)
		}
	}()
	cycleCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	err := w.monitor.Handle(cycleCtx, w.concurrency)
	if err != nil {
		logger.ErrorContext(
			cycleCtx,
			"failed to process monitor worker",
			"monitorWorker.go",
			"runCycle",
			err,
		)
		return
	}
	logger.InfoContext(
		cycleCtx,
		"monitor worker completed successfully",
		"monitorWorker.go",
		"runCycle",
	)
}
