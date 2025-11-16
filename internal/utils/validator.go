// internal/utils/validator.go
package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func NewValidator() *validator.Validate {
	return validator.New()
}

var ValidatorSet = wire.NewSet(
	NewValidator,
)