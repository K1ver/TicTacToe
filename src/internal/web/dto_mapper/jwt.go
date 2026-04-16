package dto_mapper

import (
	"app/internal/domain/model"
	"app/internal/web/dto"
)

func JwtRequestToDto(domain *model.JwtRequest) *dto.JwtRequest {
	return &dto.JwtRequest{
		Login:    domain.Login,
		Password: domain.Password,
	}
}

func JwtRequestToDomainFromDto(dto *dto.JwtRequest) *model.JwtRequest {
	return &model.JwtRequest{
		Login:    dto.Login,
		Password: dto.Password,
	}
}

func JwtResponseToDto(domain *model.JwtResponse) *dto.JwtResponse {
	return &dto.JwtResponse{
		AccessToken:  domain.AccessToken,
		RefreshToken: domain.RefreshToken,
	}
}

func JwtResponseToDomainFromDto(dto *dto.JwtResponse) *model.JwtResponse {
	return &model.JwtResponse{
		AccessToken:  dto.AccessToken,
		RefreshToken: dto.RefreshToken,
	}
}

func RefreshJwtRequestToDto(domain *model.RefreshJwtRequest) *dto.RefreshJwtRequest {
	return &dto.RefreshJwtRequest{RefreshToken: domain.RefreshToken}
}

func RefreshJwtRequestToDomainFromDto(dto *dto.RefreshJwtRequest) *model.RefreshJwtRequest {
	return &model.RefreshJwtRequest{RefreshToken: dto.RefreshToken}
}
