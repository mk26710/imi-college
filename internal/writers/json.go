package writers

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	err := enc.Encode(data)

	if err != nil {
		return err
	}

	return nil
}
