package models

type DtoPlayState struct {
	AdventureId    string      `json:"adventureId"`
	Player         string      `json:"player"`
	ActualPosition DtoPosition `json:"actualPosition"`
}
