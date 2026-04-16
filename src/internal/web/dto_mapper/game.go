package dto_mapper

import (
	"app/internal/domain/model"
	"app/internal/web/dto"
)

func GameToDto(domain *model.Game) *dto.Game {
	r, c := domain.Field.GetSize()
	f := dto.NewGameField(r, c)

	for i := range r {
		for j := range c {
			f[i][j] = domain.Field[i][j]
		}
	}

	return &dto.Game{
		Id:          domain.Id,
		Field:       f,
		Status:      dto.GameStatus(domain.Status),
		PlayerX:     domain.PlayerX,
		PlayerO:     domain.PlayerO,
		CurrentTurn: domain.CurrentTurn,
		WinnerId:    domain.WinnerId,
		CreatedAt:   domain.CreatedAt,
	}
}

func GameToDomainFromDto(dto *dto.Game) *model.Game {
	return &model.Game{
		Id:          dto.Id,
		Field:       model.GameField(dto.Field),
		Status:      model.GameStatus(dto.Status),
		PlayerX:     dto.PlayerX,
		PlayerO:     dto.PlayerO,
		CurrentTurn: dto.CurrentTurn,
		WinnerId:    dto.WinnerId,
		CreatedAt:   dto.CreatedAt,
	}
}
