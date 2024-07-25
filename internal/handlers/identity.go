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
	"gorm.io/gorm"
)

type IdentityDocsHanlder struct {
	db *gorm.DB
}

// GET /users/{userId}/documents/identity
func (h *IdentityDocsHanlder) Read(w http.ResponseWriter, r *http.Request) error {
	_, targetUser, err := httpx.GetUsersFromPathWithUAC(h.db, r, "userId", permissions.PermissionViewUser)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpx.NotFound()
		}
		return err
	}

	var docs []models.IdentityDoc

	if err := h.db.Where(models.IdentityDoc{UserID: targetUser.ID}).Find(&docs).Error; err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, docs)
}

type CreateIdentityBody struct {
	TypeID        int       `json:"typeId" validate:"required"`
	Series        string    `json:"series" validate:"required,gte=2"`
	Number        string    `json:"number" validate:"required,gte=2"`
	Issuer        string    `json:"issuer" validate:"required,gte=2"`
	IssuedAt      date.Date `json:"issuedAt" validate:"required"`
	DivisionCode  string    `json:"divisionCode" validate:"required,gte=2"`
	NationalityID int       `json:"nationalityId" validate:"required"`
}

func (h *IdentityDocsHanlder) Create(w http.ResponseWriter, r *http.Request) error {
	if !checks.IsJson(r) {
		return httpx.BadRequest("JSON body required")
	}

	_, targetUser, err := httpx.GetUsersFromPathWithUAC(h.db, r, "userId", permissions.PermissionEditUser)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpx.NotFound()
		}
		return err
	}

	var body CreateIdentityBody

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

	newIdentity := models.IdentityDoc{
		UserID:        targetUser.ID,
		TypeID:        body.TypeID,
		Series:        body.Series,
		Number:        body.Number,
		Issuer:        body.Issuer,
		IssuedAt:      body.IssuedAt,
		DivisionCode:  body.DivisionCode,
		NationalityID: body.NationalityID,
	}

	if err := h.db.Create(&newIdentity).Error; err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, newIdentity)
}
