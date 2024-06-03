package handlers

import (
	"imi/college/internal/models"
	"imi/college/internal/writer"
	"net/http"

	"gorm.io/gorm"
)

type DictionariesHandler struct {
	db *gorm.DB
}

// GET /dictionaries/towntypes
func (h *DictionariesHandler) ReadTownTypes(w http.ResponseWriter, r *http.Request) error {
	var data []models.DictTownType

	if err := h.db.Model(&models.DictTownType{}).Find(&data).Error; err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, data)
}

func (h *DictionariesHandler) ReadRegions(w http.ResponseWriter, r *http.Request) error {
	var data []models.DictRegion

	if err := h.db.Model(models.DictRegion{}).Find(&data).Error; err != nil {
		return err
	}

	return writer.JSON(w, http.StatusOK, data)
}
