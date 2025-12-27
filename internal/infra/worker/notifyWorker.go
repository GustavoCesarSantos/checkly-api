package worker

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/nativeSQL"
	urls "GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/handlers"
	"GustavoCesarSantos/checkly-api/internal/shared/configs"
	"GustavoCesarSantos/checkly-api/internal/shared/mailer"
)

type NotifyWorker struct {
	interval time.Duration
	concurrency int
	notifyCustomer *urls.NotifyCustomer
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
	urlRepo := db.NewUrlRepository(sqlDB)
	return &NotifyWorker{
		interval: 1 * time.Minute,
		concurrency: concurrency,
		notifyCustomer: urls.NewNotifyCustomer(
			application.NewFetchPendingAlerts(alertRepo),
			application.NewMarkSent(alertRepo),
			application.NewSendEmail(m),
			application.NewScheduleRetryAlert(alertRepo),
			urlRepo,
		),
	}
}

func (w *NotifyWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	slog.Info("[NOTIFY WORKER] NotifyWorker Started", "interval", w.interval.String(), "concurrency", w.concurrency)
	for {
		select {
		case <-ctx.Done():
			slog.Info("[NOTIFY WORKER] NotifyWorker Stopped")
			return
		case <-ticker.C:
			slog.Info("[NOTIFY WORKER] NotifyWorker Tick")
			w.safeProcess(ctx)
		}
	}
}

func (w *NotifyWorker) safeProcess(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("[NOTIFY WORKER] Recovered from panic", "panic", r)
		}
	}()
	err := w.notifyCustomer.Handle(ctx, w.concurrency)
	if err != nil {
		slog.Error("[NOTIFY WORKER] Process error", "error", err)
	} else {
		slog.Info("[NOTIFY WORKER] Cycle completed successfully")
	}
}