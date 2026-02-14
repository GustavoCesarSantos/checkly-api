package urls

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
	"GustavoCesarSantos/checkly-api/internal/shared/logger"
	"GustavoCesarSantos/checkly-api/internal/shared/utils"
)

type MonitorUrls struct {
	checkUrl            *application.CheckUrl
	evaluateUrl         *application.EvaluateUrl
	fetchUrls           *application.FetchUrls
	scheduleNextCheck   *application.ScheduleNextCheck
	updateUrl           *application.UpdateUrl
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
		checkUrl:            checkUrl,
		evaluateUrl:         evaluateUrl,
		fetchUrls:           fetchUrls,
		scheduleNextCheck:   scheduleNextCheck,
		updateUrl:           updateUrl,
		updateUrlWithOutbox: updateUrlWithOutbox,
	}
}

func (m *MonitorUrls) Handle(ctx context.Context, concurrency int) error {
	opCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	now := time.Now()
	urls, err := m.fetchUrls.Execute(opCtx, now)
	if err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			logger.InfoContext(
				opCtx,
				"No pending urls found",
				"monitor-worker",
				"monitor_urls.Handle",
			)
			return nil
		}
		logger.ErrorContext(
			opCtx,
			"Failed to fetch pending urls",
			"monitor-worker",
			"monitor_urls.Handle",
			err,
		)
		return err
	}
	logger.InfoContext(
		opCtx,
		"Starting URL monitoring",
		"monitor-worker",
		"monitor_urls.Handle",
		"count", len(urls),
	)
	g, ctx := errgroup.WithContext(opCtx)
	g.SetLimit(concurrency)
	for i := range urls {
		u := &urls[i]
		g.Go(func() error {
			urlCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			result, checkErr := m.checkUrl.Execute(urlCtx, u.Address)
			if checkErr != nil {
				logger.WarnContext(
					urlCtx,
					"Check URL failed",
					"monitor-worker",
					"monitor_urls.Handle",
					"error", checkErr,
					"url", u.Address,
				)
				return nil
			}
			evaluateErr := m.evaluateUrl.Execute(u, result.IsSuccess)
			if evaluateErr != nil {
				logger.ErrorContext(
					urlCtx,
					"Failed to evaluate url",
					"monitor-worker",
					"monitor_urls.Handle",
					evaluateErr,
					"url entity", u,
				)
				return nil
			}
			scheduleErr := m.scheduleNextCheck.Execute(u, now)
			if scheduleErr != nil {
				logger.ErrorContext(
					urlCtx,
					"Failed to schedule next check",
					"monitor-worker",
					"monitor_urls.Handle",
					scheduleErr,
					"url entity", u,
				)
				return nil
			}
			if u.WentDownNow {
				updateErr := m.updateUrlWithOutbox.Execute(urlCtx, *u, dtos.UpdateUrlRequest{
					NextCheck:      u.NextCheck,
					RetryCount:     &u.RetryCount,
					DownCount:      &u.DownCount,
					StabilityCount: &u.StabilityCount,
					Status:         &u.Status,
				})
				if updateErr != nil {
					logger.ErrorContext(
						urlCtx,
						"Failed to update url with outbox",
						"monitor-worker",
						"monitor_urls.Handle",
						updateErr,
						"url entity", u,
					)
				}
				return nil
			}
			updateErr := m.updateUrl.Execute(urlCtx, u.ID, dtos.UpdateUrlRequest{
				NextCheck:      u.NextCheck,
				RetryCount:     &u.RetryCount,
				DownCount:      &u.DownCount,
				StabilityCount: &u.StabilityCount,
				Status:         &u.Status,
			})
			if updateErr != nil {
				logger.ErrorContext(
					urlCtx,
					"Failed to update url",
					"monitor-worker",
					"monitor_urls.Handle",
					updateErr,
					"url entity", u,
				)
			}
			return nil
		})
	}
	return g.Wait()
}
