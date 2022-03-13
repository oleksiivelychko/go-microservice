package helpers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

type ValidationError struct {
	validator.FieldError
}

type ValidationErrors []ValidationError

type Validation struct {
	validate *validator.Validate
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s'. Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

func (v ValidationErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func NewValidation() *Validation {
	validate := validator.New()
	_ = validate.RegisterValidation("sku", validateSKU)
	return &Validation{validate}
}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i)
	if errs == nil {
		return ValidationErrors{}
	}

	validatorErrs := errs.(validator.ValidationErrors)
	if len(validatorErrs) == 0 {
		return nil
	}

	var fmtErrs []ValidationError
	for _, err := range validatorErrs {
		ve := ValidationError{err.(validator.FieldError)}
		fmtErrs = append(fmtErrs, ve)
	}

	return fmtErrs
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[0-9]+-[0-9]+-[0-9]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	if len(matches) == 1 {
		return true
	}

	return false
}
