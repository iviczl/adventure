package models

import (
	"text-adventure/models/dbmodels"
)

type DtoPlay struct {
	Id             string `json:"id"`
	AdventureTitle string `json:"adventureTitle"`
	CreatedAt      int64  `json:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt"`
}

func PlayToDtoPlay(play dbmodels.Play) DtoPlay {
	return DtoPlay{
		Id:             play.Id.String(),
		AdventureTitle: play.AdventureTitle,
		CreatedAt:      play.CreatedAt.UnixMilli(),
		UpdatedAt:      play.UpdatedAt.UnixMilli(),
	}
}
