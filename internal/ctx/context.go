package ctx

import (
	"errors"
	"imi/college/internal/models"
	"net/http"
)

var ErrUserTokenNotFound error = errors.New("request had no user token in context")
var ErrUserDataUnattached error = errors.New("user data was not attached to token")

func GetToken(r *http.Request) (models.UserToken, error) {
	token, ok := r.Context().Value(TokenKey).(models.UserToken)
	if !ok {
		return models.UserToken{}, ErrUserTokenNotFound
	}

	return token, nil
}

func GetCurrentUser(r *http.Request) (models.User, error) {
	token, err := GetToken(r)
	if err != nil {
		return models.User{}, err
	}

	if token.UserID != token.User.ID {
		return models.User{}, ErrUserDataUnattached
	}

	return token.User, nil
}
