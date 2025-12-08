package urls

import (
	"context"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/application"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
)

type MonitorUrls struct {
	checkUrl      *application.CheckUrl
	evaluateUrl   *application.EvaluateUrl
	updateUrl     *application.UpdateUrl
	urlRepository db.IUrlRepository
}

func NewMonitorUrls(
	checkUrl *application.CheckUrl,
	evaluateUrl *application.EvaluateUrl,
	updateUrl *application.UpdateUrl,
	urlRepository db.IUrlRepository,
) *MonitorUrls {
	return &MonitorUrls{
		checkUrl:      checkUrl,
		evaluateUrl:   evaluateUrl,
		updateUrl:     updateUrl,
		urlRepository: urlRepository,
	}
}

func (mu *MonitorUrls) Handle(ctx context.Context, concurrency int) error {
	urls, err := mu.urlRepository.FindAllByNextCheck(ctx, time.Now())
	if err != nil {
		return err
	}
	if len(urls) == 0 {
		slog.Info("[worker] Nenhuma url para verificar")
		return nil
	}
	slog.Info("[worker] Urls para verificar", "count", len(urls))
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(concurrency)
	for i := range urls {
		u := &urls[i]	
		g.Go(func() error {
			result, checkErr := mu.checkUrl.Execute(u.Address)
			if checkErr != nil {
				slog.Warn("check failed", "url", u.Address, "error", checkErr)
				return checkErr
			}
			mu.evaluateUrl.Execute(u, result.IsSuccess)
			updateErr := mu.updateUrl.Execute(ctx, u.ID, dtos.UpdateUrlRequest{
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
