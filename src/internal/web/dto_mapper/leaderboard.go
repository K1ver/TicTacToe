package dto_mapper

import (
	"app/internal/domain/model"
	"app/internal/web/dto"
)

func LeaderBoardToDto(domain *model.Leaderboard) *dto.Leaderboard {
	return &dto.Leaderboard{
		UserId:  domain.UserId,
		Login:   domain.Login,
		Winrate: domain.Winrate,
	}
}

func LeaderBoardToDomainFromDto(dto *dto.Leaderboard) *model.Leaderboard {
	return &model.Leaderboard{
		UserId:  dto.UserId,
		Login:   dto.Login,
		Winrate: dto.Winrate,
	}
}
