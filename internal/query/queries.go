package query

import (
	"imi/college/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ReadUserByID(db *gorm.DB, dest *models.User, id uuid.UUID) error {
	// todo: not sure about joining address, might be a bad idea
	return db.Where(&models.User{ID: id}).Joins("Details").Joins("Address").First(dest).Error
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
