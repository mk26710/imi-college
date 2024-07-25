package handlers

import (
	"encoding/json"
	"errors"
	"imi/college/internal/checks"
	"imi/college/internal/httpx"
	"imi/college/internal/models"
	"imi/college/internal/permissions"
	"imi/college/internal/types/date"
	"imi/college/internal/validation"
	"imi/college/internal/writer"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

type CreateUserBody struct {
	FirstName  string    `json:"firstName" validate:"required,gte=2"`
	MiddleName string    `json:"middleName" validate:"required,gte=2"`
	LastName   *string   `json:"lastName" validate:"omitnil,gte=2"`
	Birthday   date.Date `json:"birthday" validate:"required"`
	GenderID   int       `json:"genderId" validate:"required"`
	UserName   string    `json:"username" validate:"required,gte=4,lte=20,username"`
	Password   string    `json:"password" validate:"required,gte=6,lte=72"`
	Email      string    `json:"email" validate:"required,email"`
	Tel        string    `json:"tel" validate:"required,e164"`
	NeedsDorm  bool      `json:"needsDorm"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) error {
	if !checks.IsJson(r) {
		return httpx.MalformedJSON()
	}

	if _, err := httpx.GetCurrentUserFromRequest(h.db, r); err == nil {
		return httpx.BadRequest("authenticated users cannot create new accounts")
	}

	var body CreateUserBody

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		return httpx.MalformedJSON()
	}

	validate := validation.NewValidator()
	if err := validate.Struct(body); err != nil {
		if cause, ok := err.(validator.ValidationErrors); ok {
			return httpx.InvalidRequest(cause)
		}
		return err
	}

	var user models.User

	txFn := func(tx *gorm.DB) error {
		user = models.User{
			Email:    body.Email,
			UserName: body.UserName,
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		details := models.UserDetails{
			UserID:     user.ID,
			FirstName:  body.FirstName,
			MiddleName: body.MiddleName,
			LastName:   body.LastName,
			GenderID:   body.GenderID,
			Birthday:   body.Birthday,
			Tel:        body.Tel,
			NeedsDorm:  body.NeedsDorm,
		}

		if err := tx.Create(&details).Error; err != nil {
			return err
		}

		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}

		password := models.Password{UserID: user.ID, Hash: string(hashedPassword)}

		if err := tx.Create(&password).Error; err != nil {
			return err
		}

		return nil
	}

	if err := h.db.Transaction(txFn); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return httpx.BadRequest("this email or username is already taken")
		}
		return err
	}

	return writer.JSON(w, http.StatusOK, user)
}

// GET /users/{id}
//
// provides information about requested user
func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) error {
	_, targetUser, err := httpx.GetUsersFromPathWithUAC(h.db, r, "userId", permissions.PermissionViewUser)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpx.NotFound()
		}
		return err
	}

	return writer.JSON(w, http.StatusOK, targetUser)
}

type UpdateUserDetailsBody struct {
	FirstName  string    `json:"firstName" validate:"required,gte=2"`
	MiddleName string    `json:"middleName" validate:"required,gte=2"`
	LastName   *string   `json:"lastName" validate:"omitnil,gte=2"`
	Birthday   date.Date `json:"birthday" validate:"required"`
	GenderID   int       `json:"genderId" validate:"required"`
	Tel        string    `json:"tel" validate:"required,e164"`
	SNILS      string    `json:"snils" validate:"required,gte=2"`
	NeedsDorm  bool      `json:"needsDorm"`
}

// PUT /users/{id}/details
func (h *UserHandler) PutDetails(w http.ResponseWriter, r *http.Request) error {
	if !checks.IsJson(r) {
		return httpx.BadRequest("request must contain JSON body")
	}

	_, targetUser, err := httpx.GetUsersFromPathWithUAC(h.db, r, "userId", permissions.PermissionEditUser)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpx.NotFound()
		}
		return err
	}

	defer r.Body.Close()

	var body UpdateUserDetailsBody

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		return httpx.BadRequest(err.Error())
	}

	validate := validation.NewValidator()
	if err := validate.Struct(body); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			return httpx.InvalidRequest(validationErr)
		}
		return err
	}

	details := models.UserDetails{
		UserID:     targetUser.ID,
		FirstName:  body.FirstName,
		MiddleName: body.MiddleName,
		LastName:   body.LastName,
		GenderID:   body.GenderID,
		Birthday:   body.Birthday,
		Tel:        body.Tel,
		SNILS:      &body.SNILS,
		NeedsDorm:  body.NeedsDorm,
	}

	// if user didn't have details associated with them whatever reason
	// db.Save() will insert struct instance from above, but if user
	// already has details we have to provide the primary key, so the
	// db.Save() would update the entry instead of creating a new one
	if targetUser.Details != nil {
		details.ID = targetUser.Details.ID
	}

	if err := h.db.Debug().Save(&details).Error; err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, map[string]any{"success": true})
}
