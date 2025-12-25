package urls

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type NotifyCustomer struct {
	fetchPendingAlerts		*application.FetchPendingAlerts
	markSent                *application.MarkSent
	sendEmail               *application.SendEmail
	scheduleRetryAlert		*application.ScheduleRetryAlert
	urlRepository			db.IUrlRepository
}

func NewNotifyCustomer(
	fetchPendingAlerts		*application.FetchPendingAlerts,
	markSent                *application.MarkSent,
	sendEmail				*application.SendEmail,
	scheduleRetryAlert		*application.ScheduleRetryAlert,
	urlRepository			db.IUrlRepository,
) *NotifyCustomer {
	return &NotifyCustomer{
		fetchPendingAlerts:		fetchPendingAlerts,
		markSent:                markSent,
		sendEmail:               sendEmail,
		scheduleRetryAlert:      scheduleRetryAlert,
		urlRepository:           urlRepository,
	}
}

func (n *NotifyCustomer) Handle(ctx context.Context, concurrency int) error {
	alerts, err := n.fetchPendingAlerts.Execute(ctx, 100)
	if err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(concurrency)
	for i := range alerts {
		a := &alerts[i]
		g.Go(func() error {
			sendErr := n.sendEmail.Execute(a.Payload)
			if sendErr != nil {
				retryErr := n.scheduleRetryAlert.Execute(ctx, a.ID)
				if retryErr != nil {
					return retryErr
				}
				return sendErr
			}
			updateErr := n.urlRepository.UpdateToNotified(ctx, a.UrlId)
			if updateErr != nil {
				slog.Error("failed to update url to notified", "urlID", a.UrlId, "error", updateErr)
				return updateErr
			}
			markErr := n.markSent.Execute(ctx, a.ID, time.Now())
			if markErr != nil {
				return markErr
			}
			return nil
		})
	}
	return g.Wait()
}