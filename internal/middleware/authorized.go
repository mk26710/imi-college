package middleware

import (
	"context"
	"imi/college/internal/contextkeys"
	"imi/college/internal/handlers"
	"imi/college/internal/models"
	"imi/college/internal/writers"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func writeError(w http.ResponseWriter) {
	writers.Json(w, http.StatusForbidden, handlers.Forbidden())
}

func EnsureUserSession(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) == 0 {
				writeError(w)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w)
				return
			}

			headerParts := strings.SplitN(authHeader, " ", 2)
			if headerParts == nil {
				writeError(w)
				return
			}

			providedToken := headerParts[1]
			if len(providedToken) < 1 {
				writeError(w)
				return
			}

			var session models.UserToken

			if err := db.Where(&models.UserToken{Token: providedToken}).Preload("User").First(&session).Error; err != nil {
				writeError(w)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.TokenKey, session)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(h)
	}
}
