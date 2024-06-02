package ctx

import (
	"errors"
	"imi/college/internal/models"
	"net/http"
)

var ErrUserNotFound error = errors.New("user data is not attached to request context")

func GetCurrentUser(r *http.Request) (models.User, error) {
	user, ok := r.Context().Value(UserKey).(models.User)
	if !ok {
		return models.User{}, ErrUserNotFound
	}

	return user, nil
}
