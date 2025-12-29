package worker

import (
	"context"
	"database/sql"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/nativeSQL"
	urls "GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/handlers"
	"GustavoCesarSantos/checkly-api/internal/shared/configs"
	"GustavoCesarSantos/checkly-api/internal/shared/logger"
	"GustavoCesarSantos/checkly-api/internal/shared/mailer"
)

type NotifyWorker struct {
	interval time.Duration
	concurrency int
	notifyCustomer *urls.NotifyCustomer
	running		chan struct{}
}

func NewNotifyWorker(sqlDB *sql.DB, concurrency int) *NotifyWorker {
	mailerConfigs := configs.LoadMailerConfig()
	m := mailer.NewMailer(
		mailerConfigs.Host,
		mailerConfigs.Port,
		mailerConfigs.Username,
		mailerConfigs.Password,
		mailerConfigs.Sender,
	)
	alertRepo := db.NewAlertOutboxRepository(sqlDB)
	return &NotifyWorker{
		interval: 1 * time.Minute,
		concurrency: concurrency,
		running: make(chan struct{}, 1),
		notifyCustomer: urls.NewNotifyCustomer(
			application.NewFetchPendingAlerts(alertRepo),
			application.NewMarkSent(alertRepo),
			application.NewSendEmail(m),
			application.NewScheduleRetryAlert(alertRepo),
		),
	}
}

func (w *NotifyWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	logger.InfoContext(
		ctx,
		"Notify worker started",
		"notifyWorker.go",
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
				"Notify worker stopped",
				"notifyWorker.go",
				"Start",
			)
			return
		case <-ticker.C:
			select {
			case w.running <- struct{}{}:
				logger.InfoContext(
					ctx,
					"Notify worker tick",
					"notifyWorker.go",
					"Start",
				)
				go w.runCycle(ctx)
			default:
				logger.WarnContext(
					ctx,
					"Previous notify worker cycle still running, skipping this tick",
					"notifyWorker.go",
					"Start",
				)
			}
		}
	}
}

func (w *NotifyWorker) waitForRunning() {
	select {
	case w.running <- struct{}{}:
		<-w.running
	default:
	}
}

func (w *NotifyWorker) runCycle(ctx context.Context) {
	defer func() {
		<-w.running
		if r := recover(); r != nil {
			logger.Warn(
				"recovered from panic",
				"notifyWorker.go",
				"runCycle",
				"panic", r,
			)
		}
	}()
	cycleCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	err := w.notifyCustomer.Handle(cycleCtx, w.concurrency)
	if err != nil {
		logger.ErrorContext(
			cycleCtx,
			"failed to process notify worker",
			"notifyWorker.go",
			"runCycle",
			err,
		)
		return
	}
	logger.InfoContext(
		ctx,
		"notify worker completed successfully",
		"notifyWorker.go",
		"runCycle",
	)
}