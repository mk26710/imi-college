package handlers

import (
	"encoding/json"
	"errors"
	"imi/college/internal/checks"
	"imi/college/internal/ctx"
	"imi/college/internal/httpx"
	"imi/college/internal/models"
	"imi/college/internal/permissions"
	"imi/college/internal/validation"
	"imi/college/internal/writer"
	"net/http"
	"time"

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
	Birthday   time.Time `json:"birthday" validate:"required"`
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
	currentUser, err := ctx.GetCurrentUser(r)
	if err != nil {
		return err
	}

	targetUser, err := httpx.GetTargetUserFromPathValue(h.db, r, "id")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpx.NotFound()
		}
		return err
	}

	if targetUser.ID != currentUser.ID {
		if !permissions.HasViewUser(currentUser.Permissions) {
			return httpx.Forbidden()
		}
	}

	return writer.JSON(w, http.StatusOK, targetUser)
}
