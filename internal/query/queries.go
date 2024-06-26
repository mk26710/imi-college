package query

import (
	"imi/college/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ReadUserByID(db *gorm.DB, dest *models.User, id uuid.UUID) error {
	return db.Where(&models.User{ID: id}).Joins("Details").First(dest).Error
}

func GetUserByID(db *gorm.DB, id uuid.UUID) (models.User, error) {
	var user models.User

	err := ReadUserByID(db, &user, id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetTokenByValue(db *gorm.DB, value string) (models.UserToken, error) {
	var token models.UserToken

	if err := db.Where(&models.UserToken{Token: value}).First(&token).Error; err != nil {
		return models.UserToken{}, err
	}

	return token, nil
}

func GetUserAddressByUserID(db *gorm.DB, userID uuid.UUID) (models.UserAddress, error) {
	var addr models.UserAddress

	if err := db.Where(&models.UserAddress{UserID: userID}).Joins("Region").Joins("TownType").First(&addr).Error; err != nil {
		return models.UserAddress{}, err
	}

	return addr, nil
}

func GetDefaultAppStatus(db *gorm.DB) (models.DictAppStatus, error) {
	var status models.DictAppStatus

	if err := db.Where(&models.DictAppStatus{IsDefault: true}).First(&status).Error; err != nil {
		return models.DictAppStatus{}, err
	}

	return status, nil
}
