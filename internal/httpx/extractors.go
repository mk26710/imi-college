package httpx

import (
	"fmt"
	"imi/college/internal/ctx"
	"imi/college/internal/models"
	"imi/college/internal/query"
	"imi/college/internal/security"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetCurrentUserFromRequest(db *gorm.DB, r *http.Request) (models.User, error) {
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

func GetTargetUserFromPathValue(db *gorm.DB, r *http.Request, param string) (models.User, error) {
	currentUser, err := ctx.GetCurrentUser(r)
	if err != nil {
		return models.User{}, err
	}

	pathValue := r.PathValue(param)

	if pathValue == "@me" {
		return currentUser, nil
	}

	id, err := uuid.Parse(pathValue)
	if err != nil {
		return models.User{}, err
	}

	if id == currentUser.ID {
		return currentUser, nil
	}

	user, err := query.GetUserByID(db, id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
