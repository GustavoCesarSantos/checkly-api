package db

import "GustavoCesarSantos/checkly-api/internal/modules/urls/domain"

type IUrlRepository interface {
	Save(url *domain.Url) error
}