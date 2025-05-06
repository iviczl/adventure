package models

import (
	"text-adventure/models/dbmodels"
	"time"
)

type DtoPlay struct {
	Id             Guid      `gorm:"primaryKey;" json:"id"`
	AdventureTitle string    `gorm:"size:60;not null" json:"adventureTitle"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func PlayToDtoPlay(play dbmodels.Play) DtoPlay {
	return DtoPlay{
		Id:             play.Id,
		AdventureTitle: play.AdventureTitle,
		CreatedAt:      play.CreatedAt,
		UpdatedAt:      play.UpdatedAt,
	}
}
