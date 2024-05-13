package handlers

import (
	"fmt"
	"imi/college/internal/contextkeys"
	"imi/college/internal/models"
	"imi/college/internal/writers"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewFilesHandler(db *gorm.DB) *FilesHandler {
	return &FilesHandler{db}
}

func IsSupportedImageType(file multipart.File) (bool, error) {
	// only 512 bytes are needed to find out if file is an image
	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return false, err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, err
	}

	switch http.DetectContentType(buf) {
	case "image/png", "image/jpeg", "image/jpg":
		return true, nil
	default:
		return false, nil
	}
}

func SaveUserImage(image multipart.File, userID uuid.UUID, filename string) error {
	baseDir := ".file-uploads"
	userDir := fmt.Sprintf("%s/%s", baseDir, userID)

	if _, err := image.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		os.Mkdir(baseDir, 0755)
	} else if err != nil {
		return err
	}

	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		os.Mkdir(userDir, 0755)
	} else if err != nil {
		return err
	}

	attachmentPath := fmt.Sprintf("%s/%s", userDir, filename)

	out, err := os.Create(attachmentPath)
	if err != nil {
		return err
	}

	defer out.Close()

	if _, err = io.Copy(out, image); err != nil {
		return err
	}

	return nil
}

type FilesHandler struct {
	db *gorm.DB
}

// This handler requires a request body to be a form
// attachment - a file user attaches
func (h *FilesHandler) CreateFile(w http.ResponseWriter, r *http.Request) error {
	token, ok := r.Context().Value(contextkeys.TokenKey).(models.UserToken)
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

	isImage, err := IsSupportedImageType(attachment)
	if err != nil {
		return err
	}

	if !isImage {
		return MalformedForm()
	}

	if err := SaveUserImage(attachment, token.UserID, handler.Filename); err != nil {
		return err
	}

	// TODO: create sensible response
	writers.Json(w, 200, map[string]any{"success": true})

	return nil
}
