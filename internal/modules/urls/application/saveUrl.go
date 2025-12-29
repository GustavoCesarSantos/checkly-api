package application

import (
	"context"
	"fmt"
	"time"

	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
)

// SaveUrl é responsável por criar e persistir
// uma nova URL monitorada no sistema.
//
// Ele define o estado inicial da URL com base
// no resultado da primeira verificação HTTP
// e agenda o primeiro NextCheck.
//
// Regras:
// - URLs saudáveis iniciam com StatusHealthy
// - URLs não saudáveis iniciam com StatusDegraded
// - NextCheck segue o intervalo configurado
// - URLs degradadas iniciam com NextCheck em 1 minuto
//
// Este serviço orquestra a criação da entidade
// e delega a persistência ao repositório.
type SaveUrl struct {
	urlRepository db.IUrlRepository
}

// NewSaveUrl cria uma nova instância do serviço SaveUrl.
func NewSaveUrl(urlRepository db.IUrlRepository) *SaveUrl {
	return &SaveUrl{
		urlRepository,
	}
}

// Execute cria uma nova URL a partir dos dados de entrada,
// aplica o estado inicial e persiste no repositório.
//
// Retorna a entidade criada ou erro em caso de falha
// na persistência.
func (s *SaveUrl) Execute(ctx context.Context, input dtos.CreateUrlRequest, isHealthy bool) (*domain.Url, error) {
	status := domain.StatusHealthy
	nextCheck := time.Now().Add(time.Duration(input.Interval) * time.Minute)
	if !isHealthy {
		status = domain.StatusDegraded
		nextCheck = time.Now().Add(1 * time.Minute)
	}
	url := domain.NewUrl(
		input.Address,
		input.Interval,
		input.RetryLimit,
		input.Contact,
		status,
		nextCheck,
	)
	err := s.urlRepository.Save(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("saveUrl: %w", err)
	}
	return url, nil
}
