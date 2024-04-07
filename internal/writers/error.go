package writers

import (
	"encoding/json"
	"net/http"
)

type errorBody struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func Error(w http.ResponseWriter, message string, statusCode int) (int, error) {
	body, err := json.Marshal(errorBody{Error: message, Code: statusCode})
	if err != nil {
		return 0, err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return w.Write(body)
}
