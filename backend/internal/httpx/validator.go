package httpx

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"portfolio/backend/internal/models"
)

// =========================
// validation error type
// =========================

type ValidationErrors map[string]string

// =========================
// validator instance
// =========================

var validate *validator.Validate

func init() {

	validate = validator.New()

	// =========================
	// json tag name support
	// =========================

	validate.RegisterTagNameFunc(
		func(fld reflect.StructField) string {

			name := strings.SplitN(
				fld.Tag.Get("json"),
				",",
				2,
			)[0]

			if name == "-" {
				return ""
			}

			return name
		},
	)

	// =========================
	// custom validators
	// =========================

	validate.RegisterValidation(
		"task_status",
		func(fl validator.FieldLevel) bool {

			status := models.TaskStatus(
				fl.Field().String(),
			)

			return status.IsValid()
		},
	)
}

// =========================
// ValidateStruct
// =========================

func ValidateStruct(
	v any,
) ValidationErrors {

	err := validate.Struct(v)

	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors

	if !errors.As(err, &ve) {

		return ValidationErrors{
			"_error": "invalid validation input",
		}
	}

	validationErrors := make(
		ValidationErrors,
	)

	for _, e := range ve {

		field := e.Field()

		switch e.Tag() {

		case "required":

			validationErrors[field] =
				"is required"

		case "email":

			validationErrors[field] =
				"must be valid email"

		case "min":

			validationErrors[field] =
				"too short"

		case "max":

			validationErrors[field] =
				"too long"

		case "oneof",
			"task_status":

			validationErrors[field] =
				"invalid value"

		default:

			validationErrors[field] =
				"invalid"
		}
	}

	return validationErrors
}
