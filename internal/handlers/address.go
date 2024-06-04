package handlers

import (
	"errors"
	"imi/college/internal/ctx"
	"imi/college/internal/query"
	"imi/college/internal/writer"
	"net/http"

	"gorm.io/gorm"
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
