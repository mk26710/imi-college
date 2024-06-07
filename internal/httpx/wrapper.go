package httpx

import (
	"imi/college/internal/writer"
	"log/slog"
	"net/http"
)

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
