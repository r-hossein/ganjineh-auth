// internal/utils/validator.go
package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func NewValidator() *validator.Validate {
	v := validator.New()

	// ثبت ولیدیتور سفارشی فارسی
	v.RegisterValidation("persian", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()

		// فقط حروف فارسی + فاصله
    	re := regexp.MustCompile(`^[\x{0600}-\x{06FF}\s]+$`)
		return re.MatchString(value)
	})

	return v
}

var ValidatorSet = wire.NewSet(
	NewValidator,
)