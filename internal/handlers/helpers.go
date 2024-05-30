package handlers

import (
	"imi/college/internal/writer"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type APIError struct {
	Success bool `json:"success"` // todo: probably better switch it, as there's no such field on successful api calls
	Status  int  `json:"status"`
	Message any  `json:"message"`
	Details any  `json:"details,omitempty"`

	cause error `json:"-"`
}

func (e APIError) Error() string {
	var msg string

	if e.cause == nil {
		switch message := e.Message.(type) {
		case string:
			msg = message
		case []string:
			msg = strings.Join(message, "\n")
		default:
			msg = ""
		}
	} else {
		msg = e.cause.Error()
	}

	return msg
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func APIHandler(h APIFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writer.JSON(w, apiErr.Status, apiErr)
			} else {
				errResp := APIError{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
				writer.JSON(w, http.StatusInternalServerError, errResp)
			}

			slog.Error("HTTP API Error", "err", err.Error(), "path", r.URL.Path)
		}
	}

	return fn
}

func UnprocessableEntity() APIError {
	return APIError{
		Status:  http.StatusUnprocessableEntity,
		Message: "Request body contains unprocessable entity",
	}
}

func MalformedForm() APIError {
	return APIError{
		Status:  http.StatusUnprocessableEntity,
		Message: "Request body contains malformed form data",
	}
}

func MalformedJSON() APIError {
	return APIError{
		Status:  http.StatusUnprocessableEntity,
		Message: "Request body contains malformed json",
	}
}

func Forbidden() APIError {
	return APIError{
		Status:  http.StatusForbidden,
		Message: "Forbidden",
	}
}

func TooLarge() APIError {
	return APIError{
		Status:  http.StatusRequestEntityTooLarge,
		Message: "Request entity is too large",
	}
}

func BadRequest(reason string) APIError {
	return APIError{
		Status:  http.StatusBadRequest,
		Message: reason,
	}
}

func NotFound() APIError {
	return APIError{
		Status:  http.StatusNotFound,
		Message: "Not found",
	}
}

func InvalidCredentials(cause error) APIError {
	return APIError{
		Status:  http.StatusUnauthorized,
		Message: "Invalid credentials",
		cause:   cause,
	}
}

func InvalidRequest(cause validator.ValidationErrors) APIError {
	type errEntry struct {
		Field     string `json:"field"`
		ActualTag string `json:"tag"`
		Param     string `json:"param,omitempty"`
	}

	m := make([]errEntry, 0)

	for _, err := range cause {
		m = append(m, errEntry{
			Field:     err.Field(),
			ActualTag: err.ActualTag(),
			Param:     err.Param(),
		})
	}

	return APIError{
		Status:  http.StatusBadRequest,
		Message: "Request body is invalid",
		Details: m,
	}
}
