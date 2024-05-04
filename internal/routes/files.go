package routes

import (
	"fmt"
	"imi/college/internal/middleware"
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

func (h *FilesHandler) CreateFile(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(middleware.SessionKey).(models.UserSession)
	if !ok {
		writers.Error(w, "Unable to obtain session data.", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	// allow 12MB uploads maximum
	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)

	attachment, handler, err := r.FormFile("attachment")
	if err != nil {
		if _, ok := err.(*http.MaxBytesError); ok {
			writers.Error(w, "Attachment exceeds file size limit.", http.StatusRequestEntityTooLarge)
			return
		}

		writers.Error(w, "Attachment could not be read.", http.StatusInternalServerError)
		return
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
		writers.Error(w, "Could not save uploaded file.", http.StatusInternalServerError)
		return
	}

	defer out.Close()

	_, err = io.Copy(out, attachment)
	if err != nil {
		writers.Error(w, "Could not copy the uploaded file.", http.StatusInternalServerError)
		return
	}

	// TODO: create sensible response
	writers.Json(w, 200, map[string]any{"success": true})
}
