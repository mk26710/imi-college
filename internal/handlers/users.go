package handlers

import (
	"encoding/json"
	"errors"
	"imi/college/internal/checks"
	"imi/college/internal/models"
	"imi/college/internal/security"
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
	UserName   string `json:"username" validate:"required,gte=4,lte=20"`
	Tel        string `json:"tel" validate:"required,gte=10,lte=15"`
	FirstName  string `json:"firstName" validate:"required"`
	MiddleName string `json:"middleName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	Password   string `json:"password" validate:"required,gte=6,lte=72"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	if !checks.IsJson(r) {
		return MalformedJSON()
	}

	var body NewUserBody

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		return MalformedJSON()
	}

	validate := validator.New()
	err := validate.Struct(body)

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		return InvalidRequest(validationErr)
	} else if err != nil {
		return err
	}

	var user models.User

	txErr := h.db.Transaction(func(tx *gorm.DB) error {
		user = models.User{
			Email:      body.Email,
			UserName:   body.UserName,
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
		return BadRequest("This email or username is already taken")
	} else if txErr != nil {
		return txErr
	}

	return writers.Json(w, http.StatusOK, user)
}

func (h *UserHandler) ReadUser(w http.ResponseWriter, r *http.Request) error {
	pathValueID := r.PathValue("id")
	if len(pathValueID) == 0 {
		return BadRequest("User ID path parameter nto found")
	}

	id, err := uuid.Parse(pathValueID)
	if err != nil {
		return BadRequest("Provided user ID is incorrect.")
	}

	var user models.User

	if err := h.db.Where(&models.User{ID: id}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFound()
		}

		return err
	}

	return writers.Json(w, http.StatusOK, user)
}

type NewSessionBody struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,gte=6,lte=72"`
}

func (h *UserHandler) CreateUserToken(w http.ResponseWriter, r *http.Request) error {
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
	err := validate.Struct(body)

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		return InvalidRequest(validationErr)
	} else if err != nil {
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

	return writers.Json(w, http.StatusOK, userToken)
}
