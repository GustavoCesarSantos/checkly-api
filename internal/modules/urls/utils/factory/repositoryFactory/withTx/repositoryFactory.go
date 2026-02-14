package factory_withtx

import (
	"database/sql"

	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
	"GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/nativeSQL"
	factory "GustavoCesarSantos/checkly-api/internal/modules/urls/utils/factory/repositoryFactory/interface"
)

type repositoryFactory struct {
	tx *sql.Tx
}

func NewRepositoryFactory(tx *sql.Tx) factory.IRepositoryFactory {
	return &repositoryFactory{tx: tx}
}

func (r *repositoryFactory) Urls() db.IUrlRepository {
	return nativeSQL.NewUrlRepository(r.tx)
}

func (r *repositoryFactory) AlertOutbox() db.IAlertOutboxRepository {
	return nativeSQL.NewAlertOutboxRepository(r.tx)
}
