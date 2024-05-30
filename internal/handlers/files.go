package handlers

import (
	"fmt"
	"imi/college/internal/ctx"
	"imi/college/internal/writer"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// checks if provided file is in image of supported type and
// will return APIError if file is unsupported or there was
// a different error while performing a mimetype check
//
// if the image is supported will return the mimetype
func ValidateImageType(file multipart.File) (string, error) {
	// 512 bytes needed as per DetectContentType docs
	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return "", err
	}

	// we should seek file reader back to start
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	mime := http.DetectContentType(buf)

	switch mime {
	case "image/png", "image/jpeg", "image/jpg":
		return mime, nil
	default:
		return "", UnprocessableEntity()
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

	attachmentPath := path.Join(userDir, filename)

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

// struct of the http handler for files
type FilesHandler struct {
	db *gorm.DB
}

func NewFilesHandler(db *gorm.DB) *FilesHandler {
	return &FilesHandler{db}
}

// This handler requires a request body to be a form
func (h *FilesHandler) CreateFile(w http.ResponseWriter, r *http.Request) error {
	user, err := ctx.GetUser(r)
	if err != nil {
		return err
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

	// check if provided file is an image of supported type
	if _, err := ValidateImageType(attachment); err != nil {
		return nil
	}

	if err := SaveUserImage(attachment, user.ID, handler.Filename); err != nil {
		return err
	}

	// TODO: create sensible response
	return writer.JSON(w, 200, map[string]any{"success": true})
}
