// internal/httpx/error.go

package httpx

import (
	"errors"
	"net/http"

	"portfolio/backend/internal/service"
)

func HandleError(
	w http.ResponseWriter,
	err error,
) {

	switch {

	// =========================
	// common
	// =========================

	case errors.Is(err, service.ErrInvalidID):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeInvalidID,
			"invalid id",
		)

	// =========================
	// user
	// =========================

	case errors.Is(err, service.ErrUserNotFound):

		WriteError(
			w,
			http.StatusNotFound,
			CodeUserNotFound,
			"user not found",
		)

	case errors.Is(err, service.ErrDuplicateEmail):

		WriteError(
			w,
			http.StatusConflict,
			CodeDuplicateEmail,
			"email already exists",
		)

	// =========================
	// task
	// =========================

	case errors.Is(err, service.ErrTaskNotFound):

		WriteError(
			w,
			http.StatusNotFound,
			CodeTaskNotFound,
			"task not found",
		)

	case errors.Is(err, service.ErrInvalidUserID):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeInvalidUserID,
			"invalid user_id",
		)

	case errors.Is(err, service.ErrInvalidStatus):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeValidationError,
			"invalid status",
		)

	case errors.Is(err, service.ErrForeignKeyViolation):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeInvalidUserID,
			"invalid user_id",
		)

	// =========================
	// fallback
	// =========================

	default:

		WriteError(
			w,
			http.StatusInternalServerError,
			CodeInternalServerError,
			"internal server error",
		)
	}
}
