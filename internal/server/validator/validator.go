package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Validate[T any](v T) []*ValidationError {
	rawErrs := validate.Struct(v)

	if rawErrs == nil {
		return nil
	}

	validationErrs := make([]*ValidationError, len(rawErrs.(validator.ValidationErrors)))

	for i, err := range rawErrs.(validator.ValidationErrors) {
		fieldName := err.Field()

		var errorMessage string
		if err.Param() != "" {
			errorMessage = fmt.Sprintf("%s %s", err.Tag(), err.Param())
		} else {
			errorMessage = err.Tag()
		}

		validationErrs[i] = &ValidationError{
			Field: fieldName,
			Error: errorMessage,
		}
	}

	return validationErrs
}
