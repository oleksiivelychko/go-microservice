package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/oleksiivelychko/go-utils/validation"
	"regexp"
)

func NewValidation() *validation.Validate {
	validate := validator.New()
	_ = validate.RegisterValidation("sku", validateSKU)
	return &validation.Validate{ValidatorValidate: validate}
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^([0-9]{3})+-([0-9]{3})+-([0-9]{3})$`)
	return re.MatchString(fl.Field().String())
}
