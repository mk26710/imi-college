package middleware

import (
	"context"
	"imi/college/internal/ctx"
	"imi/college/internal/enum"
	"imi/college/internal/handlers"
	"imi/college/internal/models"
	"imi/college/internal/security"
	"imi/college/internal/writer"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func writeError(w http.ResponseWriter) {
	data := handlers.Unauthorized()
	writer.JSON(w, data.Status, data)
}

func RequireUser(db *gorm.DB) func(next http.Handler) http.Handler {
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

func writeBadPermissions(w http.ResponseWriter) {
	data := handlers.Forbidden()
	writer.JSON(w, data.Status, data)
}

func RequirePermissions(required int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user, err := ctx.GetCurrentUser(r)
			if err != nil {
				writeBadPermissions(w)
				return
			}

			if enum.HasPermissions(user.Permissions, enum.PermissionAdmin) {
				next.ServeHTTP(w, r)
				return
			}

			if enum.HasPermissions(user.Permissions, required) {
				next.ServeHTTP(w, r)
				return
			}

			writeBadPermissions(w)
		}

		return http.HandlerFunc(fn)
	}
}
