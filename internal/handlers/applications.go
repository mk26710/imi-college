package handlers

import (
	"encoding/json"
	"imi/college/internal/httpx"
	"imi/college/internal/models"
	"imi/college/internal/permissions"
	"imi/college/internal/query"
	"imi/college/internal/validation"
	"imi/college/internal/writer"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationsHandler struct {
	db *gorm.DB
}

// GET /users/{userId}/applications
func (h *ApplicationsHandler) Read(w http.ResponseWriter, r *http.Request) error {
	_, targetUser, err := httpx.GetUsersFromPathWithUAC(h.db, r, "userId", permissions.PermissionViewUser)
	if err != nil {
		return err
	}

	var apps []models.Application

	if err := h.db.Where(&models.Application{UserID: targetUser.ID}).Find(&apps).Error; err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, apps)
}

type CreateApplicationBody struct {
	MajorID    uuid.UUID `json:"majorId" validate:"required"`
	EduLevelID int       `json:"eduLevelId" validate:"required"`
}

func (h *ApplicationsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	_, targetUser, err := httpx.GetUsersFromPathWithUAC(h.db, r, "userId", permissions.PermissionEditUser)
	if err != nil {
		return err
	}

	var body CreateApplicationBody

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

	status, err := query.GetDefaultAppStatus(h.db)
	if err != nil {
		return err
	}

	application := models.Application{
		UserID:     targetUser.ID,
		MajorID:    body.MajorID,
		EduLevelID: body.EduLevelID,
		StatusID:   status.ID,
	}

	if err := h.db.Create(&application).Error; err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, application)
}
