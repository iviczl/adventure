package models

type Npc struct {
	Code                    string       `json:"code"`
	Name                    string       `json:"name"`
	Description             string       `json:"description"`
	State                   string       `json:"state"`
	LoseState               string       `json:"loseState"`
	Interacted              bool         `json:"interacted"`
	PositionCode            string       `json:"positionCode"`
	Items                   []*Item      `json:"items"`
	Attributes              []*Attribute `json:"attributes"`
	Abilities               []*Ability   `json:"abilities"`
	EnteringActions         []*Action    `json:"enteringActions"`
	LeavingActions          []*Action    `json:"leavingActions"`
	AvailableActions        []*Action    `json:"availableActions"`
	PlayerEnteringReactions []*Action    `json:"playerEnteringReactions"`
	PlayerReactions         []*Action    `json:"playerReactions"`
}

func (npc *Npc) GetAction(actionCode string) *Action {
	for i := range npc.EnteringActions {
		if npc.EnteringActions[i].Code == actionCode {
			return npc.EnteringActions[i]
		}
	}
	for i := range npc.LeavingActions {
		if npc.LeavingActions[i].Code == actionCode {
			return npc.LeavingActions[i]
		}
	}
	for i := range npc.AvailableActions {
		if npc.AvailableActions[i].Code == actionCode {
			return npc.AvailableActions[i]
		}
	}
	for i := range npc.PlayerEnteringReactions {
		if npc.PlayerEnteringReactions[i].Code == actionCode {
			return npc.PlayerEnteringReactions[i]
		}
	}
	for i := range npc.PlayerReactions {
		if npc.PlayerReactions[i].Code == actionCode {
			return npc.PlayerReactions[i]
		}
	}
	return nil
}

func (npc *Npc) GetVitalAttribute() *Attribute {
	for _, attribute := range npc.Attributes {
		if attribute.Vital {
			return attribute
		}
	}
	return nil
}

func (npc *Npc) GetAttribute(attributeName string) *Attribute {
	for i := range npc.Attributes {
		if npc.Attributes[i].Name == attributeName {
			return npc.Attributes[i]
		}
	}
	return nil
}
