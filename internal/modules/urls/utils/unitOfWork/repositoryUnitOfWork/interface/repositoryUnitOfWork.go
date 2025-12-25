package unitOfWork

import (
	"context"

	factory "GustavoCesarSantos/checkly-api/internal/modules/urls/utils/factory/repositoryFactory/interface"
)

type IRepositoryUnitOfWork interface {
	WithTx(ctx context.Context, fn func(f factory.IRepositoryFactory) error) error
}