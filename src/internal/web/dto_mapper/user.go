package dto_mapper

import (
	"app/internal/domain/model"
	"app/internal/web/dto"
)

func UserToDto(domain *model.User) *dto.User {
	return &dto.User{
		Id:       domain.Id,
		Login:    domain.Login,
		Password: domain.Password,
	}
}

func UserToDomainFromDto(dto *dto.User) *model.User {
	return &model.User{
		Id:       dto.Id,
		Login:    dto.Login,
		Password: dto.Password,
	}
}
