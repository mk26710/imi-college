package query

import (
	"imi/college/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetUserByID(db *gorm.DB, dest *models.User, id uuid.UUID) error {
	return db.Where(&models.User{ID: id}).First(dest).Error
}

func GetUserIdentityById(db *gorm.DB, id uuid.UUID) (models.UserIdentity, error) {
	var identity models.UserIdentity

	err := db.Where(&models.UserIdentity{UserID: id}).First(&identity).Error

	return identity, err
}
