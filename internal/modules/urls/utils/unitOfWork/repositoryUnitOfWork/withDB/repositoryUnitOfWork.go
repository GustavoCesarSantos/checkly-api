package unitOfWork_withDB

import (
	"context"
	"database/sql"

	factory "GustavoCesarSantos/checkly-api/internal/modules/urls/utils/factory/repositoryFactory/interface"
	factory_withtx "GustavoCesarSantos/checkly-api/internal/modules/urls/utils/factory/repositoryFactory/withTx"
	unitOfWork "GustavoCesarSantos/checkly-api/internal/modules/urls/utils/unitOfWork/repositoryUnitOfWork/interface"
)

type repositoryUnitOfWork struct {
	db *sql.DB
}

func NewRepositoryFactory(db *sql.DB) unitOfWork.IRepositoryUnitOfWork {
	return &repositoryUnitOfWork{db: db}
}

func (r *repositoryUnitOfWork) WithTx(ctx context.Context, fn func(f factory.IRepositoryFactory) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	txFactory := factory_withtx.NewRepositoryFactory(tx)
	if err = fn(txFactory); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
