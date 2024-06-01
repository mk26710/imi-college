package query

import (
	"imi/college/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetUserByID(db *gorm.DB, dest *models.User, id uuid.UUID) error {
	return db.Where(&models.User{ID: id}).First(dest).Error
}
