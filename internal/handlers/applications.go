package handlers

import (
	"encoding/json"
	"errors"
	"imi/college/internal/httpx"
	"imi/college/internal/models"
	"imi/college/internal/permissions"
	"imi/college/internal/query"
	"imi/college/internal/validation"
	"imi/college/internal/writer"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	application := models.Application{
		UserID:     targetUser.ID,
		MajorID:    body.MajorID,
		EduLevelID: body.EduLevelID,
	}

	txFn := func(tx *gorm.DB) error {
		status, err := query.GetDefaultAppStatus(tx)
		if err != nil {
			return err
		}

		application.StatusID = status.ID

		var topPriorityApp models.Application

		if err := tx.Where(&models.Application{UserID: targetUser.ID}).Order("priority DESC").Limit(1).Find(&topPriorityApp).Error; err != nil {
			return err
		}

		application.Priority = topPriorityApp.Priority + 1

		if err := h.db.Create(&application).Error; err != nil {
			return err
		}

		return nil
	}

	if err := h.db.Transaction(txFn); err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, application)
}

// DELETE /users/{userId}/applications/{appId}
func (h *ApplicationsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	_, targetUser, err := httpx.GetUsersFromPathWithUAC(h.db, r, "userId", permissions.PermissionEditUser)
	if err != nil {
		return err
	}

	appId, err := uuid.Parse(chi.URLParam(r, "appId"))
	if err != nil {
		return httpx.UnprocessableEntity()
	}

	var targetApp models.Application

	txFn := func(tx *gorm.DB) error {
		if err := tx.Where(&models.Application{UserID: targetUser.ID, ID: appId}).First(&targetApp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return httpx.NotFound()
			}
			return err
		}

		if err := tx.Delete(&targetApp).Error; err != nil {
			return err
		}

		return tx.
			Model(&models.Application{}).
			Where(&models.Application{UserID: targetUser.ID}).
			Where(gorm.Expr("priority > ?", targetApp.Priority)).
			UpdateColumn("priority", gorm.Expr("priority - 1")).
			Error
	}

	if err := h.db.Transaction(txFn); err != nil {
		return err
	}

	if err := h.db.Transaction(txFn); err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, targetApp)
}
