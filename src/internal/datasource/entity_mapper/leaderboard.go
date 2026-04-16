package entity_mapper

import (
	"app/internal/datasource/entity"
	"app/internal/domain/model"
)

func LeaderBoardToEntity(domain *model.Leaderboard) *entity.Leaderboard {
	return &entity.Leaderboard{
		UserId:  domain.UserId,
		Login:   domain.Login,
		Winrate: domain.Winrate,
	}
}

func LeaderBoardToDomainFromEntity(entity *entity.Leaderboard) *model.Leaderboard {
	return &model.Leaderboard{
		UserId:  entity.UserId,
		Login:   entity.Login,
		Winrate: entity.Winrate,
	}
}
