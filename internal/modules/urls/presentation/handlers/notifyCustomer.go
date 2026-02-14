package urls

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	"GustavoCesarSantos/checkly-api/internal/shared/logger"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type NotifyCustomer struct {
	fetchPendingAlerts *application.FetchPendingAlerts
	markSent           *application.MarkSent
	sendEmail          *application.SendEmail
	scheduleRetryAlert *application.ScheduleRetryAlert
}

func NewNotifyCustomer(
	fetchPendingAlerts *application.FetchPendingAlerts,
	markSent *application.MarkSent,
	sendEmail *application.SendEmail,
	scheduleRetryAlert *application.ScheduleRetryAlert,
) *NotifyCustomer {
	return &NotifyCustomer{
		fetchPendingAlerts: fetchPendingAlerts,
		markSent:           markSent,
		sendEmail:          sendEmail,
		scheduleRetryAlert: scheduleRetryAlert,
	}
}

func (n *NotifyCustomer) Handle(ctx context.Context, concurrency int) error {
	alerts, err := n.fetchPendingAlerts.Execute(ctx, 100)
	if err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			logger.InfoContext(
				ctx,
				"No pending alerts found",
				"notify-worker",
				"notify_customer.Handle",
			)
			return nil
		}
		logger.ErrorContext(
			ctx,
			"Failed to fetch pending alerts",
			"notify-worker",
			"notify_customer.Handle",
			err,
		)
		return err
	}
	logger.InfoContext(
		ctx,
		"Starting URL monitoring",
		"notify-worker",
		"notify_customer.Handle",
		"count", len(alerts),
	)
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(concurrency)
	for i := range alerts {
		a := &alerts[i]
		g.Go(func() error {
			sendErr := n.sendEmail.Execute(a.Payload)
			if sendErr != nil {
				logger.WarnContext(
					ctx,
					"Failed to send alert email",
					"notify-worker",
					"notify_customer.Handle",
					"error", sendErr,
					"to", a.Payload.Email,
				)
				retryErr := n.scheduleRetryAlert.Execute(ctx, a.ID)
				if retryErr != nil {
					if errors.Is(retryErr, utils.ErrRecordNotFound) {
						logger.WarnContext(
							ctx,
							"Alert not found for retry scheduling",
							"notify-worker",
							"notify_customer.Handle",
							"error", retryErr,
							"alertID", a.ID,
						)
					}
					logger.ErrorContext(
						ctx,
						"Failed to schedule retry for alert",
						"notify-worker",
						"notify_customer.Handle",
						retryErr,
						"alertID", a.ID,
					)
					return retryErr
				}
				logger.InfoContext(
					ctx,
					"Scheduled retry for alert",
					"notify-worker",
					"notify_customer.Handle",
					"alertID", a.ID,
				)
				return sendErr
			}
			markErr := n.markSent.Execute(ctx, a.ID, time.Now())
			if markErr != nil {
				if errors.Is(markErr, utils.ErrRecordNotFound) {
					logger.WarnContext(
						ctx,
						"Alert not found for marking as sent",
						"notify-worker",
						"notify_customer.Handle",
						"error", markErr,
						"alertID", a.ID,
					)
				}
				logger.ErrorContext(
					ctx,
					"Failed to mark alert as sent",
					"notify-worker",
					"notify_customer.Handle",
					markErr,
					"alertID", a.ID,
				)
				return markErr
			}
			return nil
		})
	}
	return g.Wait()
}
