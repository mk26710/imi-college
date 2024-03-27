package routes

import (
	"encoding/json"
	"errors"
	"imi/college/internal/models"
	"imi/college/internal/responses"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
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
	if r.Header.Get("Content-Type") != "application/json" {
		responses.Error(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var body NewUserBody
	var unmarshallErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&body)
	if err != nil {
		if errors.As(err, &unmarshallErr) {
			responses.Error(w, "Bad Request. Wrong Type provided for field "+unmarshallErr.Field, http.StatusBadRequest)
		} else {
			responses.Error(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	validate := validator.New()
	validationErr := validate.Struct(body)
	if validationErr != nil {
		http.Error(w, "Error during validation", http.StatusBadRequest)
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
		http.Error(w, "This email is already taken", http.StatusBadRequest)
		return
	} else if txErr != nil {
		http.Error(w, "Something went wrong, try again later", http.StatusInternalServerError)
		return
	}

	responseJson, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (h *UserHandler) ReadUser(w http.ResponseWriter, r *http.Request) {
	pathValueID := r.PathValue("id")
	if len(pathValueID) == 0 {
		http.Error(w, "User ID path parameter is not found.", http.StatusInternalServerError)
		return
	}

	id, parseErr := uuid.Parse(pathValueID)
	if parseErr != nil {
		http.Error(w, "Make sure you have request a correct user ID.", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.db.Where(&models.User{ID: id}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}

		http.Error(w, "Error while reading requested user.", http.StatusInternalServerError)
		return
	}

	resJson, marshallErr := json.Marshal(user)
	if marshallErr != nil {
		http.Error(w, "Error while converting user object to json.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
