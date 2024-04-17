package routes

import (
	"encoding/json"
	"errors"
	"imi/college/internal/checks"
	"imi/college/internal/models"
	"imi/college/internal/security"
	"imi/college/internal/writers"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SessionHandler struct {
	db *gorm.DB
}

func NewSessionHandler(db *gorm.DB) *SessionHandler {
	return &SessionHandler{db: db}
}

type NewSessionBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,lte=72"`
}

type NewSessionResponse struct {
	User    models.User        `json:"user"`
	Session models.UserSession `json:"session"`
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if !checks.IsJson(r) {
		writers.Error(w, "This endpoint requires JSON content!", http.StatusUnsupportedMediaType)
		return
	}

	var body NewSessionBody

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		writers.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(body); err != nil {
		writers.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	var user models.User

	if err := h.db.Where(&models.User{Email: body.Email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writers.Error(w, "User not found.", http.StatusNotFound)
			return
		}

		writers.Error(w, "Unexpected error while searching for user.", http.StatusInternalServerError)
		return
	}

	var password models.Password

	if err := h.db.Where(&models.Password{User: user}).First(&password).Error; err != nil {
		writers.Error(w, "Unexpected error while searching for user info.", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password.Hash), []byte(body.Password)); err != nil {
		writers.Error(w, "Wrong email or password", http.StatusUnauthorized)
		return
	}

	newToken, err := security.NewToken(security.DEFAULT_TOKEN_SIZE)
	if err != nil {
		writers.Error(w, "Unexpected errro while autheticating", http.StatusInternalServerError)
		return
	}

	userSession := models.UserSession{User: user, Token: newToken}

	if err := h.db.Create(&userSession).Error; err != nil {
		writers.Error(w, "Could not create a new session", http.StatusInternalServerError)
		return
	}

	response := NewSessionResponse{
		User:    user,
		Session: userSession,
	}

	writers.Json(w, http.StatusOK, response)
}
