package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdventurePhase string

const (
	NOT_LOADED AdventurePhase = "-1"
	STARTED    AdventurePhase = "1"
	ENDED      AdventurePhase = "2"
)

var AdventurePhases = map[string]AdventurePhase{
	"NOT_LOADED": NOT_LOADED,
	"STARTED":    STARTED,
	"ENDED":      ENDED,
}

type Adventure struct {
	// Id                uuid.UUID      `gorm:"type:uuid;primaryKey;default:lower(hex(randomblob(16)))" json:"id"`
	Id                Guid           `gorm:"primaryKey;" json:"id"`
	Code              string         `gorm:"size:10;not null" json:"code"`
	Title             string         `gorm:"size:60;not null" json:"title"`
	Description       string         `gorm:"size:100;not null" json:"description"`
	StartPositionCode string         `gorm:"size:10;not null" json:"startPositionCode"`
	Phase             AdventurePhase `gorm:"type:varchar(255);default:'-1'" json:"phase"`
	ActualPositionId  Guid           `gorm:"index" json:"actualPositionId"` // Foreign key for Position
	ActualPosition    *Position      `gorm:"foreignKey:ActualPositionId"  json:"actualPosition"`
	PlayerId          Guid           `gorm:"index" json:"playerId"` // Foreign key for Player
	Player            *Player        `gorm:"foreignKey:PlayerId" json:"player"`
	Positions         []*Position    `gorm:"foreignKey:AdventureId" json:"positions"`                        // One-to-many relation to Position
	AvailableActions  []*Action      `gorm:"foreignKey:AvailableActionsAdventureId" json:"availableActions"` // One-to-many relation to Action
}

func (a *Adventure) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	a.Id = Guid(id)
	return err
}
