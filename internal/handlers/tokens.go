package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"imi/college/internal/checks"
	"imi/college/internal/env"
	"imi/college/internal/models"
	"imi/college/internal/security"
	"imi/college/internal/writer"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TokensHandler struct {
	db *gorm.DB
}

type NewSessionBody struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,gte=6,lte=72"`
}

func NewTokensHandler(db *gorm.DB) *TokensHandler {
	return &TokensHandler{db}
}

// POST /tokens
//
// allows guests to authenticate their requests
func (h *TokensHandler) Create(w http.ResponseWriter, r *http.Request) error {
	if !checks.IsJson(r) {
		return MalformedJSON()
	}

	var body NewSessionBody

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(body); err != nil {
		if cause, casted := err.(validator.ValidationErrors); casted {
			return InvalidRequest(cause)
		}
		return err
	}

	var user models.User

	if err := h.db.Where(&models.User{UserName: body.UserName}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return InvalidCredentials(err)
		}
		return err
	}

	var password models.Password

	if err := h.db.Where(&models.Password{User: user}).First(&password).Error; err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password.Hash), []byte(body.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return InvalidCredentials(err)
		}
		return err
	}

	newToken, err := security.NewToken(security.DEFAULT_TOKEN_SIZE)
	if err != nil {
		return err
	}

	userToken := models.UserToken{User: user, Token: newToken}

	if err := h.db.Create(&userToken).Error; err != nil {
		return err
	}

	// check if requested token needs to be attached to a cookie
	if r.URL.Query().Get("cookie") == "true" {
		cookie := http.Cookie{
			Name:     "token",
			Value:    fmt.Sprintf("Bearer %v", userToken.Token),
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
			MaxAge:   int(time.Until(userToken.ExpiresAt).Seconds()),
		}

		// in production assume SSL
		if env.IsProduction() {
			cookie.Secure = true
		}

		http.SetCookie(w, &cookie)
	}

	return writer.JSON(w, http.StatusOK, userToken)
}

// DELETE /tokens
//
// allows users to delete authrentication tokens (log out process)
func (h *TokensHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	inputToken, err := security.ExtractToken(r)
	if err != nil {
		if errors.Is(err, security.ErrTokenNotFound) {
			return BadRequest("authentication token not attached to the request")
		}
		return err
	}

	txFn := func(tx *gorm.DB) error {
		var userToken models.UserToken

		if err := tx.Where(&models.UserToken{Token: inputToken}).First(&userToken).Error; err != nil {
			return err
		}

		if err := tx.Delete(userToken).Error; err != nil {
			return err
		}

		return nil
	}

	if err := h.db.Transaction(txFn); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return BadRequest("invalid token")
		}
		return err
	}

	response := map[string]any{"deleted": true}

	return writer.JSON(w, http.StatusOK, response)
}
