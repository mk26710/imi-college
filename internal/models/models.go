package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleRegular string = "regular"
	RoleStaff   string = "staff"
	RoleAdmin   string = "admin"
)

type User struct {
	ID         uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid();not null;" json:"id"`
	UserName   string     `gorm:"uniqueIndex;not null;" json:"username"`
	Email      string     `gorm:"uniqueIndex;not null;" json:"email"`
	IsVerified bool       `gorm:"default:false;not null;" json:"isVerified"`
	Role       string     `gorm:"default:'regular';not null;" json:"role"`
	Tel        *string    `json:"tel"`
	FirstName  string     `gorm:"not null;" json:"firstName"`
	MiddleName string     `gorm:"not null;" json:"middleName"`
	LastName   string     `gorm:"not null;" json:"lastName"`
	Birthday   *time.Time `json:"birthday"`
	CreatedAt  time.Time  `gorm:"default:now();not null;" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"default:now();not null;" json:"updatedAt"`
}

type Password struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid();not null;"`
	User   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID uuid.UUID `gorm:"uniqueIndex;not null;"`
	Hash   string    `gorm:"not null;"`
}

type UserSession struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid();not null;" json:"id"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null;" json:"-"`
	UserID    uuid.UUID `gorm:"not null;" json:"userId"`
	ExpiresAt time.Time `gorm:"default:now() + interval '2 days';not null;" json:"expiresAt"`
	CreatedAt time.Time `gorm:"default:now();not null;" json:"createdAt"`
	Token     string    `gorm:"uniqueIndex;not null;" json:"token"`
}
