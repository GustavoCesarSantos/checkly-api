package application

import "GustavoCesarSantos/checkly-api/internal/modules/urls/domain"

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
func (e *EvaluateUrl) Execute(url *domain.Url, httpOK bool) {
	switch {
	case url.Status == domain.StatusHealthy && !httpOK:
		url.Status = domain.StatusDegraded

	case url.Status == domain.StatusDegraded && httpOK:
		url.Status = domain.StatusRecovering
		url.RetryCount = 0

	case url.Status == domain.StatusRecovering && httpOK:
		url.StabilityCount++
		if url.StabilityCount >= 3 {
			url.Status = domain.StatusHealthy
			url.StabilityCount = 0
		}

	case url.Status == domain.StatusRecovering && !httpOK:
		url.Status = domain.StatusDegraded
		url.StabilityCount = 0

	case url.Status == domain.StatusDegraded && !httpOK:
		if url.RetryCount >= url.RetryLimit {
			url.Status = domain.StatusDown
			url.RetryCount = 0
		} else {
			url.RetryCount++
		}
	}
}
