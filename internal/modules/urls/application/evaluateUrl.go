package application

import (
	"GustavoCesarSantos/checkly-api/internal/modules/urls/domain"
	"errors"
	"fmt"
)

// EvaluateUrl é responsável por aplicar a máquina de estados
// de uma URL monitorada com base no resultado do último check HTTP.
//
// Este serviço avalia transições de estado como:
// - Healthy -> Degraded
// - Degraded -> Recovering
// - Recovering -> Healthy
// - Recovering -> Degraded
// - Degraded -> Down
//
// Além do estado, ele também atualiza:
// - RetryCount
// - StabilityCount
// - NextCheck
//
// Este serviço NÃO executa chamadas HTTP
// e NÃO persiste alterações em banco.
type EvaluateUrl struct{}

// NewEvaluateUrl cria uma nova instância do serviço EvaluateUrl.
func NewEvaluateUrl() *EvaluateUrl {
	return &EvaluateUrl{}
}

// Execute avalia o estado atual da URL e o resultado do check HTTP,
// aplicando as regras de transição de estado.
//
// Parâmetros:
// - url: entidade de domínio que será modificada
// - httpOK: indica se o último check HTTP foi bem-sucedido
//
// Observações:
// - A função modifica a entidade recebida (efeito colateral esperado)
// - O cálculo do NextCheck depende do estado final
func (e *EvaluateUrl) Execute(url *domain.Url, httpOK bool) error {
	switch {
	case url.Status == domain.StatusHealthy && httpOK:
		return nil
	case url.Status == domain.StatusHealthy && !httpOK:
		url.Status = domain.StatusDegraded
		url.RetryCount = 1
		return nil
	case url.Status == domain.StatusDegraded && httpOK:
		url.Status = domain.StatusRecovering
		url.RetryCount = 0
		url.StabilityCount = 1
		return nil
	case url.Status == domain.StatusDegraded && !httpOK:
		if url.RetryCount >= url.RetryLimit {
			url.Status = domain.StatusDown
			url.RetryCount = 0
			url.DownCount = 1
			url.WentDownNow = true
		} else {
			url.RetryCount++
		}
		return nil
	case url.Status == domain.StatusRecovering && httpOK:
		if url.StabilityCount >= 3 {
			url.Status = domain.StatusHealthy
			url.StabilityCount = 0
		} else {
			url.StabilityCount++
		}
		return nil
	case url.Status == domain.StatusRecovering && !httpOK:
		url.Status = domain.StatusDegraded
		url.StabilityCount = 0
		url.RetryCount = 1
		return nil
	case url.Status == domain.StatusDown && httpOK:
		url.Status = domain.StatusRecovering
		url.DownCount = 0
		url.StabilityCount = 1
		return nil
	case url.Status == domain.StatusDown && !httpOK:
		url.DownCount++
		return nil
	default:
		return fmt.Errorf("evaluateUrl: %w", errors.New("unhandled url status transition"))
	}
}
