package entity_mapper

import (
	"app/internal/datasource/entity"
	"app/internal/domain/model"
)

func UserToEntity(domain *model.User) *entity.User {
	return &entity.User{
		Id:       domain.Id,
		Login:    domain.Login,
		Password: domain.Password,
	}
}

func UserToDomainFromEntity(entity *entity.User) *model.User {
	return &model.User{
		Id:       entity.Id,
		Login:    entity.Login,
		Password: entity.Password,
	}
}
