package models

import (
	"encoding/json"
	"text-adventure/types"
	// "github.com/google/uuid"
	// "gorm.io/gorm"
)

type Action struct {
	Code                string                   `gorm:"size:10;not null" json:"code"`
	Description         string                   `gorm:"size:100;not null" json:"description"`
	Operation           types.ActionOperation    `gorm:"size:10;not null" json:"operation"`
	Attributes          map[string]int           `json:"attributes"`
	Abilities           map[string]bool          `json:"abilities"`
	PositionCode        string                   `gorm:"size:10;" json:"positionCode"`
	PositionDescription string                   `gorm:"size:500;" json:"positionDescription"`
	PositionVisited     bool                     `gorm:"default:true" json:"positionVisited"`
	ItemCode            string                   `gorm:"size:10" json:"itemCode"`
	ItemDescription     string                   `gorm:"size:500" json:"itemDescription"`
	ActionCode          string                   `gorm:"size:10" json:"actionCode"`
	ActionCodes         []string                 `json:"actionCodes"`
	ElseActionCodes     []string                 `json:"elseActionCodes"`
	ActionVisible       *bool                    `gorm:"default:true" json:"actionVisible"`
	ActionActive        *bool                    `gorm:"default:true" json:"actionActive"`
	NpcCode             string                   `json:"npcCode"`
	NpcInteracted       bool                     `json:"npcInteracted"`
	Value               string                   `gorm:"size:40" json:"value"`
	Function            map[string]interface{}   `gorm:"-" json:"function"`
	Functions           []map[string]interface{} `gorm:"-" json:"functions"`
	Active              *bool                    `gorm:"default:true;not null" json:"active"`
	Visible             *bool                    `gorm:"default:true;not null" json:"visible"`
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

// func (a *Action) MarshalBinary() (_ []byte, err error) {
// 	var buf bytes.Buffer
// 	enc := gob.NewEncoder(&buf)
// 	enc.Encode(p.Name)
// 	if p.Q == nil {
// 			return buf.Bytes(), nil
// 	}
// 	isCyclic := p.Q != nil && p.Q.P == p
// 	enc.Encode(isCyclic)
// 	if isCyclic {
// 			p.Q.P = nil
// 			err = enc.Encode(p.Q)
// 			p.Q.P = p
// 	} else {
// 			err = enc.Encode(p.Q)
// 	}
// 	//buf.Encode
// 	return buf.Bytes(), err
// }

// func (a *Action) UnmarshalBinary(data []byte) (err error) {
// 	buffer := bytes.NewBuffer(data)
// 	dec := gob.NewDecoder(buffer)
// 	type Alias Action
// 	var ac Alias

// 	if err = dec.Decode(&ac); err != nil {
// 		return err
// 	}
// 	*a = Action(ac)
// 	if a.Active == nil {
// 		a.Active = new(bool)
// 	}
// 	if a.Visible == nil {
// 		a.Visible = new(bool)
// 	}
// 	if a.ActionActive == nil {
// 		a.ActionActive = new(bool)
// 	}
// 	if a.ActionVisible == nil {
// 		a.ActionVisible = new(bool)
// 	}
// 	return nil
// }

func (a *Action) AdjustAction() {
	if a.Active == nil {
		a.Active = new(bool)
	}
	if a.Visible == nil {
		a.Visible = new(bool)
	}
	if a.ActionActive == nil {
		a.ActionActive = new(bool)
	}
	if a.ActionVisible == nil {
		a.ActionVisible = new(bool)
	}
}
