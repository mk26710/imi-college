package middleware

import (
	"context"
	"imi/college/internal/ctx"
	"imi/college/internal/handlers"
	"imi/college/internal/models"
	"imi/college/internal/security"
	"imi/college/internal/writers"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func writeError(w http.ResponseWriter) {
	writers.Json(w, http.StatusForbidden, handlers.Forbidden())
}

func EnsureUserSession(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			rawToken, err := security.ExtractToken(r)
			if err != nil {
				writeError(w)
				return
			}

			var userToken models.UserToken

			if err := db.Where(&models.UserToken{Token: rawToken}).Preload("User").First(&userToken).Error; err != nil {
				writeError(w)
				return
			}

			// check if found token has expired
			if userToken.ExpiresAt.Before(time.Now()) {
				db.Delete(userToken)
				writeError(w)
				return
			}

			ctx := context.WithValue(r.Context(), ctx.TokenKey, userToken)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
