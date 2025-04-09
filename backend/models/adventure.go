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

// Build constructs an Adventure object from the given data.
// func Build(data map[string]interface{}, adventure *Adventure) (*Adventure, error) {
// 	if adventure == nil {
// 		adventure = &Adventure{}
// 	}

// 	adventure.Code = data["code"].(string)
// 	adventure.Player = Player{
// 		Name: data["player"].(map[string]interface{})["name"].(string),
// 	}
// 	adventure.Title = data["title"].(string)
// 	adventure.Description = data["description"].(string)

// 	if items, ok := data["player"].(map[string]interface{})["items"]; ok {
// 		for _, itemData := range items.([]interface{}) {
// 			item := Item{}
// 			// Assuming a helper function `PrepareItem` exists to map item data
// 			PrepareItem(itemData.(map[string]interface{}), &item)
// 			adventure.Player.Items = append(adventure.Player.Items, item)
// 		}
// 	}

// 	adventure.StartPositionCode = data["start_position_code"].(string)

// 	if actions, ok := data["available_actions"]; ok {
// 		for _, actionData := range actions.([]interface{}) {
// 			action := Action{}
// 			// Assuming a helper function `PrepareAction` exists to map action data
// 			PrepareAction(actionData.(map[string]interface{}), &action)
// 			adventure.AvailableActions = append(adventure.AvailableActions, action)
// 		}
// 	}

// 	for _, positionData := range data["positions"].([]interface{}) {
// 		position := Position{}
// 		posMap := positionData.(map[string]interface{})
// 		position.Code = posMap["code"].(string)
// 		position.Description = posMap["description"].(string)
// 		position.Visited = false

// 		if endPosition, ok := posMap["end_position"]; ok {
// 			position.EndPosition = endPosition.(bool)
// 		}

// 		if items, ok := posMap["items"]; ok {
// 			for _, itemData := range items.([]interface{}) {
// 				item := Item{}
// 				PrepareItem(itemData.(map[string]interface{}), &item)
// 				position.Items = append(position.Items, item)
// 			}
// 		}

// 		if actions, ok := posMap["available_actions"]; ok {
// 			for _, actionData := range actions.([]interface{}) {
// 				action := Action{}
// 				PrepareAction(actionData.(map[string]interface{}), &action)
// 				position.AvailableActions = append(position.AvailableActions, action)
// 			}
// 		}

// 		if enteringActions, ok := posMap["entering_actions"]; ok {
// 			for _, actionData := range enteringActions.([]interface{}) {
// 				action := Action{}
// 				PrepareAction(actionData.(map[string]interface{}), &action)
// 				position.EnteringActions = append(position.EnteringActions, action)
// 			}
// 		}

// 		if leavingActions, ok := posMap["leaving_actions"]; ok {
// 			for _, actionData := range leavingActions.([]interface{}) {
// 				action := Action{}
// 				PrepareAction(actionData.(map[string]interface{}), &action)
// 				position.LeavingActions = append(position.LeavingActions, action)
// 			}
// 		}

// 		adventure.Positions = append(adventure.Positions, position)
// 	}

// 	return adventure, nil
// }

// Start initializes the Adventure.
// func (a *Adventure) Start() {
// 	action := Action{
// 		Operation:    CHANGE_POSITION, // Assuming CHANGE_POSITION is a predefined constant
// 		PositionCode: a.StartPositionCode,
// 		Active:       true,
// 	}
// 	a.Phase = STARTED
// 	action.Execute(a) // Assuming Execute is a method of Action
// }

// Do performs an action based on the given action code.
// func (a *Adventure) Do(actionCode string) error {
// 	var action *Action

// 	if a.ActualPosition != nil {
// 		action = GetActionFromPosition(a.ActualPosition, actionCode) // Assuming GetActionFromPosition is a helper function
// 	}

// 	if action == nil && len(a.AvailableActions) > 0 {
// 		for _, availableAction := range a.AvailableActions {
// 			if availableAction.Code == actionCode {
// 				action = &availableAction
// 				break
// 			}
// 		}
// 	}

// 	if action == nil {
// 		return errors.New("invalid action code in position: " + a.ActualPosition.Code + ":" + actionCode)
// 	}

// 	action.Execute(a) // Assuming Execute is a method of Action
// 	return nil
// }
