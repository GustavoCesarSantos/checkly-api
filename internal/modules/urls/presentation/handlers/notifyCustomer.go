package urls

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type NotifyCustomer struct {
	fetchPendingAlerts		*application.FetchPendingAlerts
	markSent                *application.MarkSent
	sendEmail               *application.SendEmail
	scheduleRetryAlert		*application.ScheduleRetryAlert
}

func NewNotifyCustomer(
	fetchPendingAlerts		*application.FetchPendingAlerts,
	markSent                *application.MarkSent,
	sendEmail				*application.SendEmail,
	scheduleRetryAlert		*application.ScheduleRetryAlert,
) *NotifyCustomer {
	return &NotifyCustomer{
		fetchPendingAlerts:		fetchPendingAlerts,
		markSent:                markSent,
		sendEmail:               sendEmail,
		scheduleRetryAlert:      scheduleRetryAlert,
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
			markErr := n.markSent.Execute(ctx, a.ID, time.Now())
			if markErr != nil {
				return markErr
			}
			return nil
		})
	}
	return g.Wait()
}