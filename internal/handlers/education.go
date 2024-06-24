package handlers

import (
	"encoding/json"
	"errors"
	"imi/college/internal/httpx"
	"imi/college/internal/models"
	"imi/college/internal/permissions"
	"imi/college/internal/validation"
	"imi/college/internal/writer"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type EducationDocsHandler struct {
	db *gorm.DB
}

// GET /users/{userId}/documents/education
func (h *EducationDocsHandler) Read(w http.ResponseWriter, r *http.Request) error {
	_, targerUser, err := httpx.GetUsersFromPathWithUAC(h.db, r, "userId", permissions.PermissionViewUser)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpx.NotFound()
		}
		return err
	}

	var docs []models.EducationDoc

	if err := h.db.Where(models.EducationDoc{UserID: targerUser.ID}).Find(&docs).Error; err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, docs)
}

type EducationDocBody struct {
	TypeID         int       `json:"typeId" validate:"required"`
	Series         string    `json:"series" validate:"required"`
	Number         string    `json:"number" validate:"required"`
	Issuer         string    `json:"issuer" validate:"required"`
	IssuedAt       time.Time `json:"issuedAt" validate:"required"`
	GradYear       int16     `json:"gradYear" validate:"required"`
	IssuerRegionID int       `json:"issuerRegionId" validate:"required"`
}

// POST /users/{userId}/documents/education
func (h *EducationDocsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	_, targetUser, err := httpx.GetUsersFromPathWithUAC(h.db, r, "userId", permissions.PermissionEditUser)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpx.NotFound()
		}
		return err
	}

	defer r.Body.Close()

	var body EducationDocBody

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

	newDoc := models.EducationDoc{
		UserID:         targetUser.ID,
		TypeID:         body.TypeID,
		Series:         body.Series,
		Number:         body.Number,
		Issuer:         body.Issuer,
		IssuedAt:       body.IssuedAt,
		GradYear:       body.GradYear,
		IssuerRegionID: body.IssuerRegionID,
	}

	if err := h.db.Create(&newDoc).Error; err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, newDoc)
}
