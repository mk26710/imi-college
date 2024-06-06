package extras

import (
	"fmt"
	"imi/college/internal/models"
	"imi/college/internal/query"
	"imi/college/internal/security"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func GetCurrentUser(db *gorm.DB, r *http.Request) (models.User, error) {
	rawToken, err := security.ExtractToken(r)
	if err != nil {
		return models.User{}, err
	}

	token, err := query.GetTokenByValue(db, rawToken)
	if err != nil {
		return models.User{}, err
	}

	if token.ExpiresAt.Before(time.Now()) {
		db.Delete(token)
		return models.User{}, fmt.Errorf("token has expired")
	}

	user, err := query.GetUserByID(db, token.UserID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
