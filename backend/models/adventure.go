package models

import (
	"text-adventure/types"

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
	Id                types.Guid     `gorm:"primaryKey;" json:"id"`
	Code              string         `gorm:"size:10;not null" json:"code"`
	Title             string         `gorm:"size:60;not null" json:"title"`
	Description       string         `gorm:"size:100;not null" json:"description"`
	StartPositionCode string         `gorm:"size:10;not null" json:"startPositionCode"`
	Phase             AdventurePhase `gorm:"type:varchar(255);default:'-1'" json:"phase"`
	ActualPositionId  types.Guid     `gorm:"index" json:"actualPositionId"` // Foreign key for Position
	ActualPosition    *Position      `gorm:"foreignKey:ActualPositionId"  json:"actualPosition"`
	PlayerId          types.Guid     `gorm:"index" json:"playerId"` // Foreign key for Player
	Player            *Player        `gorm:"foreignKey:PlayerId" json:"player"`
	Positions         []*Position    `gorm:"foreignKey:AdventureId" json:"positions"`                        // One-to-many relation to Position
	AvailableActions  []*Action      `gorm:"foreignKey:AvailableActionsAdventureId" json:"availableActions"` // One-to-many relation to Action
}

func (a *Adventure) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	a.Id = types.Guid(id)
	return err
}

// GetAction retrieves an action by its code from the adventure.
func (adventure *Adventure) GetAction(actionCode string) *Action {
	// Try to get the action from the positions
	action := getActionFromPositionList(adventure.Positions, actionCode)
	if action == nil {
		// If not found, try to get it from the available actions
		for i := range adventure.AvailableActions {
			if adventure.AvailableActions[i].Code == actionCode {
				return adventure.AvailableActions[i]
			}
		}
	}
	return action
}

// getActionFromPositionList retrieves an action by its code from a list of positions.
func getActionFromPositionList(positions []*Position, actionCode string) *Action {
	for i := range positions {
		action := positions[i].GetAction(actionCode)
		if action != nil {
			return action
		}
	}
	return nil
}

func (adventure *Adventure) GetPosition(positionCode string) *Position {
	for i := range adventure.Positions {
		if adventure.Positions[i].Code == positionCode {
			return adventure.Positions[i]
		}
	}
	return nil
}

// getItem retrieves an item by its code from the adventure.
func (adventure *Adventure) GetItem(itemCode string) *Item {
	// Check the player's items
	item := adventure.Player.GetItem(itemCode)
	if item != nil {
		return item
	}
	// Check the items in the current position
	if adventure.ActualPosition != nil {
		item := adventure.ActualPosition.GetItem(itemCode)
		if item != nil {
			return item
		}
	}
	// Check the items in all positions
	for i := range adventure.Positions {
		item := adventure.Positions[i].GetItem(itemCode)
		if item != nil {
			return item
		}
	}
	return nil
}
