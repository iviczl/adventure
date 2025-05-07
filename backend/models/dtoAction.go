package models

type DtoAction struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func ActionToDtoAction(a *Action) DtoAction {
	return DtoAction{
		Code:        a.Code,
		Description: a.Description,
	}
}
