package dbmodels

import (
	"time"

	"text-adventure/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id               models.Guid `gorm:"primaryKey;" json:"id"`
	UserName         string      `gorm:"size:100;not null" json:"name"`
	Email            string      `gorm:"size:100;not null" json:"email"`
	Password         string      `gorm:"size:20;not null" json:"password"`
	CreatedAt        time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	RegistrationCode string      `gorm:"size:10" json:"registrationCode"`
}

func (a *User) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	a.Id = models.Guid(id)
	return err
}
