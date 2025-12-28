package urls

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
	"GustavoCesarSantos/checkly-api/internal/shared/logger"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type MonitorUrls struct {
	checkUrl      *application.CheckUrl
	evaluateUrl   *application.EvaluateUrl
	fetchUrls *application.FetchUrls
	scheduleNextCheck *application.ScheduleNextCheck
	updateUrl     *application.UpdateUrl
	updateUrlWithOutbox *application.UpdateUrlWithOutbox
}

func NewMonitorUrls(
	checkUrl *application.CheckUrl,
	evaluateUrl *application.EvaluateUrl,
	fetchUrls *application.FetchUrls,
	scheduleNextCheck *application.ScheduleNextCheck,
	updateUrl *application.UpdateUrl,
	updateUrlWithOutbox *application.UpdateUrlWithOutbox,
) *MonitorUrls {
	return &MonitorUrls{
		checkUrl:      checkUrl,
		evaluateUrl:   evaluateUrl,
		fetchUrls:     fetchUrls,
		scheduleNextCheck: scheduleNextCheck,
		updateUrl:     updateUrl,
		updateUrlWithOutbox: updateUrlWithOutbox,
	}
}

func (m *MonitorUrls) Handle(ctx context.Context, concurrency int) error {
	urls, err := m.fetchUrls.Execute(ctx, time.Now())
	if err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			logger.InfoContext(
				ctx,
				"No pending urls found",
				"monitor-worker",
				"monitor_urls.Handle",
			)
			return nil
		}
		logger.ErrorContext(
			ctx,
			"Failed to fetch pending urls",
			"monitor-worker",
			"monitor_urls.Handle",
			err,
		)
		return err
	}
	logger.InfoContext(
		ctx,
		"Starting URL monitoring",
		"monitor-worker",
		"monitor_urls.Handle",
		"count", len(urls),
	)
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(concurrency)
	for i := range urls {
		u := &urls[i]
		g.Go(func() error {
			result, checkErr := m.checkUrl.Execute(u.Address)
			if checkErr != nil {
				logger.WarnContext(
					ctx,
					"Check URL failed",
					"monitor-worker",
					"monitor_urls.Handle",
					"url", u.Address,
					"error", checkErr,
				)
				return checkErr
			}
			evaluateErr := m.evaluateUrl.Execute(u, result.IsSuccess)
			if evaluateErr != nil {
				logger.ErrorContext(
					ctx,
					"Failed to evaluate url",
					"monitor-worker",
					"monitor_urls.Handle",
					evaluateErr,
					"url entity", u,
				)
				return evaluateErr
			}
			scheduleErr := m.scheduleNextCheck.Execute(u, time.Now())
			if scheduleErr != nil {
				logger.ErrorContext(
					ctx,
					"Failed to schedule next check",
					"monitor-worker",
					"monitor_urls.Handle",
					scheduleErr,
					"url entity", u,
				)
				return scheduleErr
			}
			if(u.Status == domain.StatusDown) {
				updateErr := m.updateUrlWithOutbox.Execute(ctx, *u, dtos.UpdateUrlRequest{
					NextCheck:      u.NextCheck,
					RetryCount:     &u.RetryCount,
					DownCount:		&u.DownCount,
					StabilityCount: &u.StabilityCount,
					Status:         &u.Status,
				})
				if updateErr != nil {
					logger.ErrorContext(
						ctx,
						"Failed to update url with outbox",
						"monitor-worker",
						"monitor_urls.Handle",
						updateErr,
						"url entity", u,
					)
					return updateErr
				}
				return nil
			}
			updateErr := m.updateUrl.Execute(ctx, u.ID, dtos.UpdateUrlRequest{
				NextCheck:      u.NextCheck,
				RetryCount:     &u.RetryCount,
				DownCount:		&u.DownCount,
				StabilityCount: &u.StabilityCount,
				Status:         &u.Status,
			})
			if updateErr != nil {
				logger.ErrorContext(
					ctx,
					"Failed to update url",
					"monitor-worker",
					"monitor_urls.Handle",
					updateErr,
					"url entity", u,
				)
				return updateErr
			}
			return nil
		})
	}
	return g.Wait()
}
