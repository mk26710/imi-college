package httpx

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Error(w http.ResponseWriter, cause error) error {
	apiError, ok := cause.(APIError)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Internal Server Error"))
		return err
	}

	response, err := json.Marshal(apiError)
	if err != nil {
		return err
	}

	w.WriteHeader(apiError.Status)
	if _, err := w.Write(response); err != nil {
		return err
	}

	return nil
}

type APIError struct {
	Status  int `json:"status"`
	Message any `json:"message"`
	Details any `json:"details,omitempty"`

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

func Unauthorized() APIError {
	return APIError{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
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
