package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt  time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"default:now()" json:"updated_at"`
	Email      string     `gorm:"uniqueIndex" json:"email"`
	IsVerified bool       `gorm:"default:false" json:"isVerified"`
	Tel        *string    `json:"tel"`
	FirstName  string     `json:"firstName"`
	MiddleName string     `json:"middleName"`
	LastName   string     `json:"lastName"`
	Birthday   *time.Time `json:"birthday"`
}

type Password struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Hash   string
}
