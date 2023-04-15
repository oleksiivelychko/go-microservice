package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

type ValidatorError struct {
	validator.FieldError
}

type ValidatorErrors []ValidatorError

type Validate struct {
	ValidatorValidate *validator.Validate
}

func (validationErr ValidatorError) Error() string {
	return fmt.Sprintf(
		"Key: '%s'. Error: field validation for '%s', failed on the '%s' tag",
		validationErr.Namespace(),
		validationErr.Field(),
		validationErr.Tag(),
	)
}

func (validationErrors ValidatorErrors) Errors() []string {
	var validationErrArray []string
	for _, err := range validationErrors {
		validationErrArray = append(validationErrArray, err.Error())
	}

	return validationErrArray
}

func (validate *Validate) Validate(i interface{}) ValidatorErrors {
	validatorValidate := validate.ValidatorValidate.Struct(i)
	if validatorValidate == nil {
		return ValidatorErrors{}
	}

	validationErrors := validatorValidate.(validator.ValidationErrors)
	if len(validationErrors) == 0 {
		return ValidatorErrors{}
	}

	var validationErrArr []ValidatorError
	for _, err := range validationErrors {
		validationError := ValidatorError{err.(validator.FieldError)}
		validationErrArr = append(validationErrArr, validationError)
	}

	return validationErrArr
}

func New() (*Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("sku", validateSKU)
	if err != nil {
		return nil, err
	}

	return &Validate{ValidatorValidate: validate}, nil
}

func validateSKU(field validator.FieldLevel) bool {
	re := regexp.MustCompile(`^([0-9]{3})+-([0-9]{3})+-([0-9]{3})$`)
	return re.MatchString(field.Field().String())
}
