package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID uint64 `json:"id" gorm:"primaryKey"`

	ExternalId  uuid.UUID `json:"external_id,omitempty"   gorm:"type:uuid"`
	FirstName   string    `json:"first_name,omitempty"    gorm:"type:varchar(100);not null"`
	LastName    string    `json:"last_name,omitempty"     gorm:"type:varchar(100);not null"`
	Email       string    `json:"email,omitempty"         gorm:"type:varchar(100);not null"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty" gorm:"type:date"`
}
