package urls

import (
	"context"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
)

type MonitorUrls struct {
	checkUrl      *application.CheckUrl
	evaluateUrl   *application.EvaluateUrl
	scheduleNextCheck *application.ScheduleNextCheck
	updateUrl     *application.UpdateUrl
	updateUrlWithOutbox *application.UpdateUrlWithOutbox
	urlRepository db.IUrlRepository
}

func NewMonitorUrls(
	checkUrl *application.CheckUrl,
	evaluateUrl *application.EvaluateUrl,
	scheduleNextCheck *application.ScheduleNextCheck,
	updateUrl *application.UpdateUrl,
	updateUrlWithOutbox *application.UpdateUrlWithOutbox,
	urlRepository db.IUrlRepository,
) *MonitorUrls {
	return &MonitorUrls{
		checkUrl:      checkUrl,
		evaluateUrl:   evaluateUrl,
		scheduleNextCheck: scheduleNextCheck,
		updateUrl:     updateUrl,
		updateUrlWithOutbox: updateUrlWithOutbox,
		urlRepository: urlRepository,
	}
}

func (m *MonitorUrls) Handle(ctx context.Context, concurrency int) error {
	urls, err := m.urlRepository.FindAllByNextCheck(ctx, time.Now())
	if err != nil {
		return err
	}
	if len(urls) == 0 {
		slog.Info("No urls to check")
		return nil
	}
	slog.Info("Urls to check", "count", len(urls))
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(concurrency)
	for i := range urls {
		u := &urls[i]
		g.Go(func() error {
			result, checkErr := m.checkUrl.Execute(u.Address)
			if checkErr != nil {
				slog.Warn("check failed", "url", u.Address, "error", checkErr)
				return checkErr
			}
			m.evaluateUrl.Execute(u, result.IsSuccess)
			m.scheduleNextCheck.Execute(u, result.IsSuccess, time.Now())
			if(u.Status == domain.StatusDown) {
				updateErr := m.updateUrlWithOutbox.Execute(ctx, *u, dtos.UpdateUrlRequest{
					NextCheck:      u.NextCheck,
					RetryCount:     &u.RetryCount,
					StabilityCount: &u.StabilityCount,
					Status:         &u.Status,
				})
				if updateErr != nil {
					slog.Error("update failed", "urlId", u.ID, "error", updateErr)
					return updateErr
				}
				return nil
			}
			updateErr := m.updateUrl.Execute(ctx, u.ID, dtos.UpdateUrlRequest{
				NextCheck:      u.NextCheck,
				RetryCount:     &u.RetryCount,
				StabilityCount: &u.StabilityCount,
				Status:         &u.Status,
			})
			if updateErr != nil {
				slog.Error("update failed", "urlId", u.ID, "error", updateErr)
				return updateErr
			}
			return nil
		})
	}
	return g.Wait()
}
