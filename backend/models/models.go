package models

// "text-adventure/types"
// "github.com/google/uuid"
// "gorm.io/gorm"

type Item struct {
	// Id          types.Guid `gorm:"primaryKey;" json:"id"`
	Code        string `gorm:"size:10;not null" json:"code"`
	Name        string `gorm:"size:40;not null" json:"name"`
	Description string `gorm:"size:200;not null" json:"description"`
	State       string `gorm:"size:40;" `
	// PlayerId    types.Guid `gorm:"index" json:"playerId"` // Foreign Key to Player
	// Player      *Player    `gorm:"foreignKey:PlayerId;" json:"player"`
	// PositionId types.Guid `gorm:"index" json:"positionId"` // Foreign Key to Position
	// Position   *Position  `gorm:"foreignKey:PositionId;" json:"position"`
}

// func (a *Item) BeforeCreate(tx *gorm.DB) error {
// 	id, err := uuid.NewRandom()
// 	a.Id = types.Guid(id)
// 	return err
// }
