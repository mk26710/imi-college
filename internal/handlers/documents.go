package handlers

import (
	"encoding/json"
	"imi/college/internal/checks"
	"imi/college/internal/ctx"
	"imi/college/internal/models"
	"imi/college/internal/validation"
	"imi/college/internal/writer"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentsHandler struct {
	db *gorm.DB
}

type CreateDocumentIdentityBody struct {
	TypeID        int       `json:"typeId" validate:"required"`
	Series        string    `json:"series" validate:"required,gte=2,lte=8"`
	Number        string    `json:"number" validate:"required"`
	Issuer        string    `json:"issuer" validate:"required"`
	IssuedAt      time.Time `json:"issuedAt" validate:"required"`
	DivisionCode  string    `json:"divisionCode" validate:"required"`
	NationalityID int       `json:"nationalityId" validate:"required"`
}

func (h *DocumentsHandler) CreateDocumentIdentity(w http.ResponseWriter, r *http.Request) error {
	if !checks.IsJson(r) {
		return MalformedJSON()
	}

	user, err := ctx.GetCurrentUser(r)
	if err != nil {
		return err
	}

	var body CreateDocumentIdentityBody

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&body); err != nil {
		return err
	}

	validate := validation.NewValidator()
	if err := validate.Struct(body); err != nil {
		if cause, ok := err.(validator.ValidationErrors); ok {
			return InvalidRequest(cause)
		}
		return err
	}

	var doc models.IdentityDoc

	txFn := func(tx *gorm.DB) error {
		doc = models.IdentityDoc{
			UserID:        user.ID,
			TypeID:        body.TypeID,
			Series:        body.Series,
			Number:        body.Number,
			Issuer:        body.Issuer,
			IssuedAt:      body.IssuedAt,
			DivisionCode:  body.DivisionCode,
			NationalityID: body.NationalityID,
			FileID:        uuid.NullUUID{Valid: false},
		}

		if err := tx.Create(&doc).Error; err != nil {
			return err
		}

		return nil
	}

	if err := h.db.Transaction(txFn); err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, doc)
}
