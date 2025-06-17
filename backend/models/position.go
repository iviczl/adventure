package models

import (
	"encoding/json"
	"slices"
	// "text-adventure/types"
	// "github.com/google/uuid"
	// "gorm.io/gorm"
)

type Position struct {
	// Id                        types.Guid `gorm:"primaryKey;" json:"id"`
	Code             string    `gorm:"size:10;not null" json:"code"`
	Description      string    `gorm:"size:500;not null" json:"description"`
	Visited          bool      `gorm:"default: false" json:"visited"`
	EndPosition      bool      `gorm:"default: false" json:"endPosition"`
	AvailableActions []*Action `gorm:"foreignKey:AvailableActionsPositionId" json:"availableActions"`
	EnteringActions  []*Action `gorm:"foreignKey:EnteringActionsPositionId" json:"enteringActions"`
	LeavingActions   []*Action `gorm:"foreignKey:LeavingActionsPositionId" json:"leavingActions"`
	Items            []*Item   `gorm:"foreignKey:PositionId" json:"items"`
	// AdventureId               types.Guid `gorm:"foreignKey:Id" json:"adventureId"`
	// ActualPositionAdventureId types.Guid `gorm:"foreignKey:Id;" json:"actualPositionAdventureId"` // Foreign Key to Adventure
	TemporaryDescription string `gorm:"size:500;" json:"temporaryDescription"`
}

// func (p *Position) BeforeCreate(tx *gorm.DB) error {
// 	id, err := uuid.NewRandom()
// 	p.Id = types.Guid(id)
// 	return err
// }

func (position *Position) GetItem(itemCode string) *Item {
	for i := range position.Items {
		if position.Items[i].Code == itemCode {
			return position.Items[i]
		}
	}
	return nil
}

func (position *Position) GetAction(actionCode string) *Action {
	// Check available actions
	for i := range position.AvailableActions {
		if position.AvailableActions[i].Code == actionCode {
			return position.AvailableActions[i]
		}
	}
	// Check entering actions
	for i := range position.EnteringActions {
		if position.EnteringActions[i].Code == actionCode {
			return position.EnteringActions[i]
		}
	}
	// Check leaving actions
	for i := range position.LeavingActions {
		if position.LeavingActions[i].Code == actionCode {
			return position.LeavingActions[i]
		}
	}
	return nil
}

func (p *Position) UnmarshalJSON(text []byte) error {
	type Alias Position
	aux := Alias{
		AvailableActions: []*Action{}, // set the default value before parsing JSON
	}
	if err := json.Unmarshal(text, &aux); err != nil {
		return err
	}
	*p = Position(aux)
	return nil
}

func AdjustedActualPosition(p *Position) *Position {
	description := p.Description
	if p.TemporaryDescription != "" {
		description = p.TemporaryDescription
		p.TemporaryDescription = ""
	}
	availableActions := slices.Collect(func(yield func(*Action) bool) {
		for i := range p.AvailableActions {
			action := p.AvailableActions[i]
			if *action.Active && action.Visible != nil && *action.Visible {
				if !yield(action) {
					return
				}
			}
		}
	})
	if availableActions == nil {
		availableActions = []*Action{}
	}

	actual := &Position{
		Code:                 p.Code,
		Description:          description,
		Visited:              p.Visited,
		EndPosition:          p.EndPosition,
		AvailableActions:     availableActions,
		EnteringActions:      p.EnteringActions,
		LeavingActions:       p.LeavingActions,
		Items:                p.Items,
		TemporaryDescription: p.TemporaryDescription,
	}
	return actual
}

// func (p *Position) MarshalJSON() ([]byte, error) {
// 	type Alias Position
// 	availableActions := slices.Collect(func(yield func(*Action) bool) {
// 		for _, a := range p.AvailableActions {
// 			if *(a.Active) && *(a.Visible) {
// 				if !yield(a) {
// 					return
// 				}
// 			}
// 		}
// 	})
// 	if availableActions == nil {
// 		availableActions = []*Action{}
// 	}
// 	return json.Marshal(&struct {
// 		AvailableActions []*Action `json:"availableActions"`
// 		*Alias
// 	}{
// 		AvailableActions: availableActions,
// 		Alias:            (*Alias)(p),
// 	})
// }
