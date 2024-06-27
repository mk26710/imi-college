package handlers

import (
	"imi/college/internal/models"
	"imi/college/internal/writer"
	"net/http"
	"time"

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

	writer.SetCacheControlSWR(w, 24*time.Hour, 6*time.Hour)
	return writer.JSON(w, http.StatusOK, data)
}

// GET /dictionaries/regions
func (h *DictionariesHandler) ReadRegions(w http.ResponseWriter, r *http.Request) error {
	var data []models.DictRegion

	if err := h.db.Model(models.DictRegion{}).Find(&data).Error; err != nil {
		return err
	}

	writer.SetCacheControlSWR(w, 24*time.Hour, 6*time.Hour)
	return writer.JSON(w, http.StatusOK, data)
}

// GET /dictionaries/genders
func (h *DictionariesHandler) ReadGenders(w http.ResponseWriter, r *http.Request) error {
	var data []models.DictGender

	if err := h.db.Model(models.DictGender{}).Find(&data).Error; err != nil {
		return err
	}

	writer.SetCacheControlSWR(w, 24*time.Hour, 6*time.Hour)
	return writer.JSON(w, http.StatusOK, data)
}

// GET /dictionaries/edulevels
func (h *DictionariesHandler) ReadEduLevels(w http.ResponseWriter, r *http.Request) error {
	var data []models.DictEduLevel

	if err := h.db.Model(models.DictEduLevel{}).Find(&data).Error; err != nil {
		return err
	}

	writer.SetCacheControlSWR(w, 24*time.Hour, 6*time.Hour)
	return writer.JSON(w, http.StatusOK, data)
}

// GET /dictionaries/majors
func (h *DictionariesHandler) ReadMajors(w http.ResponseWriter, r *http.Request) error {
	var data []models.CollegeMajor

	if err := h.db.Model(models.CollegeMajor{}).Find(&data).Error; err != nil {
		return err
	}

	writer.SetCacheControlSWR(w, 24*time.Hour, 6*time.Hour)
	return writer.JSON(w, http.StatusOK, data)
}

func (h *DictionariesHandler) ReadAppStatuses(w http.ResponseWriter, r *http.Request) error {
	var data []models.DictAppStatus

	if err := h.db.Model(models.DictAppStatus{}).Find(&data).Error; err != nil {
		return err
	}

	writer.SetCacheControlSWR(w, 24*time.Hour, 6*time.Hour)
	return writer.JSON(w, http.StatusOK, data)
}
