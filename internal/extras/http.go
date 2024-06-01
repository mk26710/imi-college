package extras

import (
	"fmt"
	"imi/college/internal/models"
	"imi/college/internal/security"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func UserFromHttp(db *gorm.DB, r *http.Request) (models.User, error) {
	rawToken, err := security.ExtractToken(r)
	if err != nil {
		return models.User{}, err
	}

	var token models.UserToken

	if err := db.Where(&models.UserToken{Token: rawToken}).First(&token).Error; err != nil {
		return models.User{}, err
	}

	if token.ExpiresAt.Before(time.Now()) {
		db.Delete(token)
		return models.User{}, fmt.Errorf("token has expired")
	}

	var user models.User

	if err := db.Where(&models.User{ID: token.UserID}).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}
