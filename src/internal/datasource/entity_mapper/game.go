package entity_mapper

import (
	"app/internal/datasource/entity"
	"app/internal/domain/model"
	"encoding/json"
)

func GameToDomainFromEntity(entity *entity.Game) *model.Game {
	var modelField model.GameField
	_ = json.Unmarshal(entity.Field, &modelField)
	return &model.Game{
		Id:          entity.ID,
		Field:       modelField,
		Status:      model.GameStatus(entity.Status),
		PlayerX:     entity.PlayerX,
		PlayerO:     entity.PlayerO,
		CurrentTurn: entity.CurrentTurn,
		WinnerId:    entity.WinnerId,
		CreatedAt:   entity.CreatedAt,
	}
}

func GameToEntity(domain *model.Game) *entity.Game {

	entityGameField, _ := json.Marshal(domain.Field)

	return &entity.Game{
		ID:          domain.Id,
		Field:       entityGameField,
		Status:      entity.GameStatus(domain.Status),
		PlayerX:     domain.PlayerX,
		PlayerO:     domain.PlayerO,
		CurrentTurn: domain.CurrentTurn,
		WinnerId:    domain.WinnerId,
		CreatedAt:   domain.CreatedAt,
	}
}
