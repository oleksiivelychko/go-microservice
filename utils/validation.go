package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/oleksiivelychko/go-utils/validator_helper"
	"regexp"
)

func NewValidation() *validator_helper.Validation {
	validate := validator.New()
	_ = validate.RegisterValidation("sku", validateSKU)
	return &validator_helper.Validation{ValidatorValidate: validate}
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^([0-9]{3})+-([0-9]{3})+-([0-9]{3})$`)
	return re.MatchString(fl.Field().String())
}
