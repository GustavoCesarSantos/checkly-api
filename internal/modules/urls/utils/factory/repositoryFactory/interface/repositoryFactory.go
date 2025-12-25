package factory

import (
	db "GustavoCesarSantos/checkly-api/internal/modules/urls/external/db/interfaces"
)

type IRepositoryFactory interface {
	Urls() db.IUrlRepository
	AlertOutbox() db.IAlertOutboxRepository
}