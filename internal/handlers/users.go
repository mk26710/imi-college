package handlers

import (
	"encoding/json"
	"errors"
	"imi/college/internal/checks"
	"imi/college/internal/ctx"
	"imi/college/internal/extras"
	"imi/college/internal/models"
	"imi/college/internal/query"
	"imi/college/internal/validation"
	"imi/college/internal/writer"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
		return MalformedJSON()
	}

	if _, err := extras.UserFromHttp(h.db, r); err == nil {
		return BadRequest("authenticated users cannot create new accounts")
	}

	var body CreateUserBody

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		return MalformedJSON()
	}

	validate := validation.NewValidator()
	if err := validate.Struct(body); err != nil {
		if cause, ok := err.(validator.ValidationErrors); ok {
			return InvalidRequest(cause)
		}
		return err
	}

	var user models.User

	txFn := func(tx *gorm.DB) error {
		user = models.User{
			Email:     body.Email,
			UserName:  body.UserName,
			NeedsDorm: body.NeedsDorm,
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
			return BadRequest("this email or username is already taken")
		}
		return err
	}

	return writer.JSON(w, http.StatusOK, user)
}

// GET /users/@me
//
// special endpoint to read user data about currently authenticated
func (h *UserHandler) ReadMe(w http.ResponseWriter, r *http.Request) error {
	user, err := ctx.GetCurrentUser(r)
	if err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, user)
}

// GET /users/{id}
//
// provides information about requested user by ID (uuid)
func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) error {
	pathValueID := r.PathValue("id")
	if len(pathValueID) == 0 {
		return BadRequest("user ID path parameter not found")
	}

	id, err := uuid.Parse(pathValueID)
	if err != nil {
		return BadRequest("provided user ID is incorrect")
	}

	user, err := query.GetUserByID(h.db, id)
	if err != nil {
		return NotFound()
	}

	return writer.JSON(w, http.StatusOK, user)
}
