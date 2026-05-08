package handlers

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {

	validate = validator.New()

	// use json tag name
	validate.RegisterTagNameFunc(func(
		fld reflect.StructField,
	) string {

		name := strings.SplitN(
			fld.Tag.Get("json"),
			",",
			2,
		)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

func ValidateStruct(v interface{}) map[string]string {

	errs := make(map[string]string)

	err := validate.Struct(v)

	if err == nil {
		return nil
	}

	for _, e := range err.(validator.ValidationErrors) {

		field := e.Field()

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
