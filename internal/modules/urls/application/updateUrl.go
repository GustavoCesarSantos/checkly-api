package application

import (
	"context"
	"fmt"

	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/presentation/dtos"
)

// UpdateUrl é responsável por aplicar atualizações
// parciais em uma URL existente.
//
// Ele recebe apenas os campos que devem ser alterados
// e delega a persistência ao repositório.
//
// Este serviço NÃO contém regras de negócio.
// Ele apenas aplica as alterações calculadas por
// outros serviços de domínio ou aplicação.
type UpdateUrl struct {
	urlRepository db.IUrlRepository
}

// NewUpdateUrl cria uma nova instância do serviço UpdateUrl.
func NewUpdateUrl(urlRepository db.IUrlRepository) *UpdateUrl {
	return &UpdateUrl{
		urlRepository,
	}
}

// Execute aplica as alterações informadas na URL identificada
// pelo urlId.
//
// Campos nil no input NÃO são atualizados.
func (u *UpdateUrl) Execute(ctx context.Context, urlId int64, input dtos.UpdateUrlRequest) error {
	params := db.UpdateUrlParams{
		NextCheck:      input.NextCheck,
		RetryCount:     input.RetryCount,
		DownCount:      input.DownCount,
		StabilityCount: input.StabilityCount,
		Status:         input.Status,
	}
	err := u.urlRepository.Update(ctx, urlId, params)
	if err != nil {
		return fmt.Errorf("updateUrl: %w", err)
	}
	return nil
}
