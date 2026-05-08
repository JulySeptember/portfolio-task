package handlers

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(v interface{}) map[string]string {

	errs := make(map[string]string)

	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	for _, e := range err.(validator.ValidationErrors) {

		field := strings.ToLower(e.Field())

		switch e.Tag() {

		case "required":
			errs[field] = "is required"

		case "email":
			errs[field] = "must be valid email"

		case "max":
			errs[field] = "too long"

		case "min":
			errs[field] = "too short"

		case "oneof":
			errs[field] = "invalid value"

		default:
			errs[field] = "invalid"
		}
	}

	return errs
}
