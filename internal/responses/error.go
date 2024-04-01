package responses

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Error(w http.ResponseWriter, message string, statusCode int) (int, error) {
	body, err := json.Marshal(ErrorBody{Message: message, Code: statusCode})
	if err != nil {
		return 0, err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return w.Write(body)
}
