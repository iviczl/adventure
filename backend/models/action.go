package models

import (
	"encoding/json"
	"text-adventure/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActionOperation string

const (
	CHANGE_POSITION                        ActionOperation = "cp"
	CHANGE_POSITION_DESCRIPTION            ActionOperation = "cpd"
	CHANGE_POSITION_TEMPORARY_DESCRIPTION  ActionOperation = "cptd"
	APPEND_POSITION_TEMPORARY_DESCRIPTION  ActionOperation = "aptd"
	PREPEND_POSITION_TEMPORARY_DESCRIPTION ActionOperation = "pptd"
	CHANGE_POSITION_VISITED                ActionOperation = "cpv"
	CHANGE_ACTION_ACTIVE                   ActionOperation = "caa"
	CHANGE_ACTION_VISIBLE                  ActionOperation = "cav"
	CHANGE_ITEM_DESCRIPTION                ActionOperation = "cad"
	CHANGE_ITEM_STATE                      ActionOperation = "cis"
	PICK_UP_ITEM                           ActionOperation = "pui"
	PUT_DOWN_ITEM                          ActionOperation = "pdi"
	MOVE_ITEM                              ActionOperation = "moi"
	ERASE_ITEM                             ActionOperation = "eri"
	CONDITIONAL                            ActionOperation = "con"
	SWITCH_CONDITIONAL                     ActionOperation = "scon"
	LIST                                   ActionOperation = "ls"
	RANDOM                                 ActionOperation = "ran"
)

type Action struct {
	Id                          types.Guid               `gorm:"primaryKey;" json:"id"`
	Code                        string                   `gorm:"size:10;not null" json:"code"`
	Description                 string                   `gorm:"size:100;not null" json:"description"`
	Operation                   ActionOperation          `gorm:"size:10;not null" json:"operation"`
	PositionCode                string                   `gorm:"size:10;" json:"positionCode"`
	PositionDescription         string                   `gorm:"size:500;" json:"positionDescription"`
	PositionVisited             bool                     `gorm:"default:true" json:"positionVisited"`
	ItemCode                    string                   `gorm:"size:10" json:"itemCode"`
	ItemDescription             string                   `gorm:"size:500" json:"itemDescription"`
	ActionCode                  string                   `gorm:"size:10" json:"actionCode"`
	ActionCodes                 []string                 `gorm:"-" json:"actionCodes"`
	ActionVisible               *bool                    `gorm:"default:true" json:"actionVisible"`
	ActionActive                *bool                    `gorm:"default:true" json:"actionActive"`
	Value                       string                   `gorm:"size:40" json:"value"`
	Function                    map[string]interface{}   `gorm:"-" json:"function"`
	Functions                   []map[string]interface{} `gorm:"-" json:"functions"`
	Active                      *bool                    `gorm:"default:true;not null" json:"active"`
	Visible                     *bool                    `gorm:"default:true;not null" json:"visible"`
	AvailableActionsPositionId  types.Guid               `gorm:"index" json:"-"` // Foreign Key to Position
	EnteringActionsPositionId   types.Guid               `gorm:"index" json:"-"` // Foreign Key to Position
	LeavingActionsPositionId    types.Guid               `gorm:"index" json:"-"` // Foreign Key to Position
	AvailableActionsAdventureId types.Guid               `gorm:"index" json:"-"` // Foreign Key to Adventure

	// Relationships
	AvailableActionsPosition  Position  `gorm:"foreignKey:AvailableActionsPositionId" json:"-"`
	EnteringActionsPosition   Position  `gorm:"foreignKey:EnteringActionsPositionId" json:"-"`
	LeavingActionsPosition    Position  `gorm:"foreignKey:LeavingActionsPositionId" json:"-"`
	AvailableActionsAdventure Adventure `gorm:"foreignKey:AvailableActionsAdventureId" json:"-"`
}

func (a *Action) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	a.Id = types.Guid(id)
	return err
}

func (a *Action) UnmarshalJSON(text []byte) error {
	type Alias Action
	activeValue := true
	visibleValue := true
	actionActiveValue := true
	actionVisibleValue := true
	aux := Alias{
		Active:        &activeValue, // set the default value before parsing JSON
		Visible:       &visibleValue,
		ActionActive:  &actionActiveValue,
		ActionVisible: &actionVisibleValue,
	}
	if err := json.Unmarshal(text, &aux); err != nil {
		return err
	}
	*a = Action(aux)
	return nil
}
