package responses

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	resp := map[string]string{"message": message}
	jsonResp, _ := json.Marshal(resp)

	w.Write(jsonResp)
}
