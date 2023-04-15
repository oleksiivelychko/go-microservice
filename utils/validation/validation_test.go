package validation

import (
	"github.com/go-playground/validator/v10"
	"testing"
)

type Test struct {
	Count int `json:"count" validate:"required,countValidator"`
}

func TestValidation_Validate(t *testing.T) {
	newValidator := validator.New()
	_ = newValidator.RegisterValidation("countValidator", func(fl validator.FieldLevel) bool {
		return fl.Field().Int() >= 0
	})

	validation := &Validate{newValidator}

	err := validation.Validate(Test{Count: 1})
	if err == nil {
		t.Fatal("Test.Count validation failed")
	}

	err = validation.Validate(Test{Count: -1})
	if err != nil {
		t.Logf("Test.Count validation failed: %s\n", err.Errors())
	} else {
		t.Fatal("Test.Count validation did not invoke")
	}
}
