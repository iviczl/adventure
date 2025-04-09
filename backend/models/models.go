package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	Id          Guid      `gorm:"primaryKey;" json:"id"`
	Code        string    `gorm:"size:10;not null" json:"code"`
	Name        string    `gorm:"size:40;not null" json:"name"`
	Description string    `gorm:"size:200;not null" json:"description"`
	State       string    `gorm:"size:40;" `
	PlayerId    Guid      `gorm:"index" json:"playerId"` // Foreign Key to Player
	Player      *Player   `gorm:"foreignKey:PlayerId;" json:"player"`
	PositionId  Guid      `gorm:"index" json:"positionId"` // Foreign Key to Position
	Position    *Position `gorm:"foreignKey:PositionId;" json:"position"`
}

func (a *Item) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	a.Id = Guid(id)
	return err
}

type Player struct {
	Id         Guid         `gorm:"primaryKey;" json:"id"`
	Name       string       `gorm:"size:40;not null" json:"name"`
	Items      []*Item      `gorm:"foreignKey:PlayerId" json:"items"`
	Adventures []*Adventure `gorm:"foreignKey:PlayerId" json:"adventures"`
}

func (p *Player) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	p.Id = Guid(id)
	return err
}
