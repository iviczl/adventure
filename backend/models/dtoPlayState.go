package models

type DtoPlayState struct {
	AdventureId    string      `json:"adventureId"`
	Player         DtoPlayer   `json:"player"`
	ActualPosition DtoPosition `json:"actualPosition"`
}
