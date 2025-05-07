package dbmodels

import (
	"text-adventure/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Play struct {
	Id             types.Guid `gorm:"primaryKey;" json:"id"`
	UserId         types.Guid `gorm:"foreignKey:Id" json:"userId"`
	AdventureTitle string     `gorm:"size:60;not null" json:"adventureTitle"`
	Adventure      []byte     `gorm:"type:text" json:"adventure"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (a *Play) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	a.Id = types.Guid(id)
	return err
}
