package models

type DtoPosition struct {
	Code             string      `json:"code"`
	Description      string      `json:"description"`
	EndPosition      bool        `json:"endPosition"`
	AvailableActions []DtoAction `json:"availableActions"`
}

func PositionToDtoPosition(p *Position) DtoPosition {
	actions := make([]DtoAction, len(p.AvailableActions))
	for i, action := range p.AvailableActions {
		actions[i] = ActionToDtoAction(action)
	}
	return DtoPosition{
		Code:             p.Code,
		Description:      p.Description,
		EndPosition:      p.EndPosition,
		AvailableActions: actions,
	}
}
