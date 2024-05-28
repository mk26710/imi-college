package middleware

import (
	"context"
	"imi/college/internal/ctx"
	"imi/college/internal/handlers"
	"imi/college/internal/models"
	"imi/college/internal/writers"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

func writeError(w http.ResponseWriter) {
	writers.Json(w, http.StatusForbidden, handlers.Forbidden())
}

func EnsureUserSession(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var inputToken string

			// attempt to read token from cookie first
			if cookie, err := r.Cookie("token"); err == nil {
				inputToken = cookie.Value
			}

			// if cookie wasn't read or empty attempt reading header
			if header := r.Header.Get("Authorization"); len(header) > 0 && len(inputToken) == 0 {
				inputToken = header
			}

			// if token is empty then there's no token
			if len(inputToken) == 0 {
				writeError(w)
				return
			}

			// make sure token has prefix and cut it
			rawToken, isCut := strings.CutPrefix(inputToken, "Bearer ")
			if !isCut {
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
