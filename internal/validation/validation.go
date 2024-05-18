package validation

import (
	"slices"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("username", ValidateUsername)

	return validate
}

func ValidateUsernameRune(r rune) bool {
	valid := []rune("AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz_1234567890")

	return slices.Contains(valid, r)
}

func ValidateUsername(fl validator.FieldLevel) bool {
	val := fl.Field().String()

	// if at least rune is not valid return fals immediately
	for _, r := range val {
		ok := ValidateUsernameRune(r)
		if !ok {
			return false
		}
	}

	return true
}
