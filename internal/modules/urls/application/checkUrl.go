package application

import (
	"context"
	"net/http"
	"time"
)

// CheckUrl é responsável por realizar uma verificação HTTP
// em uma URL externa.
//
// Ele executa uma requisição GET com timeout fixo e retorna
// o status HTTP obtido, além de indicar se a verificação
// foi considerada bem-sucedida.
//
// Regras:
// - Timeout máximo de 3 segundos
// - Status codes >= 400 são considerados falha
// - Erros de rede são propagados
//
// Este serviço NÃO avalia estado de domínio (Healthy, Degraded, etc).
// Ele apenas executa a verificação técnica da URL.
type CheckUrl struct{}

// NewCheckUrl cria uma nova instância do serviço CheckUrl.
func NewCheckUrl() *CheckUrl {
	return &CheckUrl{}
}

// CheckUrlResult representa o resultado técnico da verificação HTTP.
type CheckUrlResult struct {
	IsSuccess  bool
	StatusCode int
}

// Execute realiza a chamada HTTP para a URL informada
// e retorna o resultado da verificação.
//
// Erros retornados indicam falhas técnicas (timeout, DNS, etc).
// Falhas HTTP (status >= 400) NÃO retornam erro, apenas IsSuccess=false.
func (c *CheckUrl) Execute(url string) (CheckUrlResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CheckUrlResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return CheckUrlResult{
			IsSuccess:  false,
			StatusCode: resp.StatusCode,
		}, nil
	}
	return CheckUrlResult{
		IsSuccess:  true,
		StatusCode: resp.StatusCode,
	}, nil
}
