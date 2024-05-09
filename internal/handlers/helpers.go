package handlers

import (
	"imi/college/internal/writers"
	"log/slog"
	"net/http"
	"strings"
)

type APIError struct {
	Status  int   `json:"status"`
	Message any   `json:"message"`
	Details []any `json:"details,omitempty"`
}

func (e APIError) Error() string {
	var msg string

	switch message := e.Message.(type) {
	case string:
		msg = message
	case []string:
		msg = strings.Join(message, "\n")
	default:
		msg = ""
	}

	return msg
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func APIHandler(h APIFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writers.Json(w, apiErr.Status, apiErr)
			} else {
				errResp := APIError{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
				writers.Json(w, http.StatusInternalServerError, errResp)
			}

			slog.Error("HTTP API Error", "err", err.Error(), "path", r.URL.Path)
		}
	}

	return fn
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
