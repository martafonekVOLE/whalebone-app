package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user's model.
type User struct {
	gorm.Model

	ID uint64 `json:"id" gorm:"primaryKey"`

	ExternalId uuid.UUID `json:"external_id,omitempty" gorm:"type:uuid"`
	// Name represents the user's full display name.
	// Splitting into FirstName/LastName is preferred for normalized storage,
	// but since I don't control the upstream data shape, a single field
	// avoids premature parsing of unknown formats.
	Name        *string   `json:"name,omitempty"          gorm:"type:varchar(200);not null"`
	Email       *string   `json:"email,omitempty"         gorm:"type:varchar(100);not null"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty" gorm:"type:date"`
}
