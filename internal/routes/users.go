package routes

import (
	"encoding/json"
	"errors"
	"imi/college/internal/checks"
	"imi/college/internal/models"
	"imi/college/internal/writers"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db}
}

type NewUserBody struct {
	Email      string `json:"email" validate:"required,email"`
	Tel        string `json:"tel" validate:"required,gte=10,lte=15"`
	FirstName  string `json:"firstName" validate:"required"`
	MiddleName string `json:"middleName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	Password   string `json:"password" validate:"required,gte=6,lte=72"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if !checks.IsJson(r) {
		writers.Error(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var body NewUserBody

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		writers.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	validationErr := validate.Struct(body)
	if validationErr != nil {
		writers.Error(w, "Error during validation", http.StatusBadRequest)
		return
	}

	var user models.User

	txErr := h.db.Transaction(func(tx *gorm.DB) error {
		user = models.User{
			Email:      body.Email,
			Tel:        &body.Tel,
			FirstName:  body.FirstName,
			MiddleName: body.MiddleName,
			LastName:   body.LastName,
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}

		if err := tx.Create(&models.Password{UserID: user.ID, Hash: string(hashedPassword)}).Error; err != nil {
			return err
		}

		return nil
	})

	if errors.Is(txErr, gorm.ErrDuplicatedKey) {
		writers.Error(w, "This email is already taken.", http.StatusBadRequest)
		return
	} else if txErr != nil {
		writers.Error(w, "Something went wrong, try again later.", http.StatusInternalServerError)
		return
	}

	writers.Json(w, http.StatusOK, user)
}

func (h *UserHandler) ReadUser(w http.ResponseWriter, r *http.Request) {
	pathValueID := r.PathValue("id")
	if len(pathValueID) == 0 {
		writers.Error(w, "User ID path parameter is not found.", http.StatusInternalServerError)
		return
	}

	id, err := uuid.Parse(pathValueID)
	if err != nil {
		writers.Error(w, "Make sure you have requested a correct user ID.", http.StatusBadRequest)
		return
	}

	var user models.User

	if err := h.db.Where(&models.User{ID: id}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writers.Error(w, "User not found.", http.StatusNotFound)
			return
		}

		writers.Error(w, "Error while reading requested user.", http.StatusInternalServerError)
		return
	}

	writers.Json(w, http.StatusOK, user)
}
