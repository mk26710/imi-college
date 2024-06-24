package httpx

import (
	"fmt"
	"imi/college/internal/ctx"
	"imi/college/internal/models"
	"imi/college/internal/permissions"
	"imi/college/internal/query"
	"imi/college/internal/security"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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

	pathValue := chi.URLParam(r, param)

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

// the function will try to get current user from request's context
// and if current user is present it will try to get target user
// from path parameter's value; once target user is acquired the
// function will perform user acceess control check via permissions
// by comparing current user's permissions with provided required
// permissions
//
// intended for use with paths like /users/{id}, /users/{id}/address and etc.
func GetUsersFromPathWithUAC(db *gorm.DB, r *http.Request, param string, required int64) (models.User, models.User, error) {
	currentUser, err := ctx.GetCurrentUser(r)
	if err != nil {
		return models.User{}, models.User{}, err
	}

	targetUser, err := GetTargetUserFromPathValue(db, r, param)
	if err != nil {
		return models.User{}, models.User{}, err
	}

	if targetUser.ID != currentUser.ID {
		if !permissions.HasPermissions(currentUser.Permissions, required) {
			return models.User{}, models.User{}, Forbidden()
		}
	}

	return currentUser, targetUser, nil
}
