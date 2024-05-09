package handlers

import (
	"fmt"
	"imi/college/internal/contextkeys"
	"imi/college/internal/models"
	"imi/college/internal/writers"
	"io"
	"net/http"
	"os"

	"gorm.io/gorm"
)

type FilesHandler struct {
	db *gorm.DB
}

func NewFilesHandler(db *gorm.DB) *FilesHandler {
	return &FilesHandler{db}
}

func (h *FilesHandler) CreateFile(w http.ResponseWriter, r *http.Request) error {
	session, ok := r.Context().Value(contextkeys.TokenKey).(models.UserToken)
	if !ok {
		return fmt.Errorf("unable to obtain session data")
	}

	defer r.Body.Close()

	// allow 12MB uploads maximum
	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)

	attachment, handler, err := r.FormFile("attachment")
	if err != nil {
		if _, ok := err.(*http.MaxBytesError); ok {
			return TooLarge()
		}

		return err
	}

	defer attachment.Close()

	baseDir := ".file-uploads"
	userDir := fmt.Sprintf("%s/%s", baseDir, session.UserID)

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		os.Mkdir(baseDir, 0755)
	}

	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		os.Mkdir(userDir, 0755)
	}

	attachmentPath := fmt.Sprintf("%s/%s", userDir, handler.Filename)

	out, err := os.Create(attachmentPath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, attachment)
	if err != nil {
		return err
	}

	// TODO: create sensible response
	writers.Json(w, 200, map[string]any{"success": true})

	return nil
}
