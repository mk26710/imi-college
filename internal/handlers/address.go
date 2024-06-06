package handlers

import (
	"encoding/json"
	"errors"
	"imi/college/internal/checks"
	"imi/college/internal/ctx"
	"imi/college/internal/extras"
	"imi/college/internal/models"
	"imi/college/internal/permissions"
	"imi/college/internal/query"
	"imi/college/internal/validation"
	"imi/college/internal/writer"
	"net/http"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AddressHandler struct {
	db *gorm.DB
}

// GET /users/@me/address
func (h *AddressHandler) ReadMe(w http.ResponseWriter, r *http.Request) error {
	user, err := ctx.GetCurrentUser(r)
	if err != nil {
		return err
	}

	addr, err := query.GetUserAddressByUserID(h.db, user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFound()
		}
		return err
	}

	return writer.JSON(w, http.StatusOK, addr)
}

type AddressBody struct {
	RegionID   int    `json:"regionID" validate:"required"`
	TownTypeID int    `json:"townTypeID" validate:"required"`
	Town       string `json:"town" validate:"required"`
	Address    string `json:"address" validate:"required"`
	PostCode   string `json:"postCode" validate:"required"`
}

func (h *AddressHandler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) error {
	if !checks.IsJson(r) {
		return BadRequest("JSON body required")
	}

	currentUser, err := ctx.GetCurrentUser(r)
	if err != nil {
		return err
	}

	targetUser, err := extras.GetTargetUserFromPathValue(h.db, r, "id")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFound()
		}
		return err
	}

	if targetUser.ID != currentUser.ID {
		if !permissions.HasEditUser(currentUser.Permissions) {
			return Forbidden()
		}
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var body AddressBody

	if err := decoder.Decode(&body); err != nil {
		return BadRequest("couldn't parse request body")
	}

	validate := validation.NewValidator()
	if err := validate.Struct(body); err != nil {
		if cause, ok := err.(validator.ValidationErrors); ok {
			return InvalidRequest(cause)
		}
		return err
	}

	newAddr := models.UserAddress{
		UserID:     targetUser.ID,
		RegionID:   body.RegionID,
		TownTypeID: body.TownTypeID,
		Town:       body.Town,
		Address:    body.Address,
		PostCode:   body.PostCode,
	}
	if err := h.db.Debug().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"region_id", "town_type_id", "town", "post_code", "address"}),
	}).Create(&newAddr).Error; err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, map[string]any{"success": true})
}
