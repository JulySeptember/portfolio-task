// internal/httpx/error.go

package httpx

import (
	"errors"
	"net/http"

	"portfolio/backend/internal/apperr"
)

func HandleError(
	w http.ResponseWriter,
	err error,
) {

	switch {

	// =========================
	// common
	// =========================

	case errors.Is(
		err,
		apperr.ErrInvalidID,
	):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeInvalidID,
			"invalid id",
		)

	case errors.Is(
		err,
		apperr.ErrValidation,
	):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeValidationError,
			"validation error",
		)

	// =========================
	// user
	// =========================

	case errors.Is(
		err,
		apperr.ErrUserNotFound,
	):

		WriteError(
			w,
			http.StatusNotFound,
			CodeUserNotFound,
			"user not found",
		)

	case errors.Is(
		err,
		apperr.ErrDuplicateEmail,
	):

		WriteError(
			w,
			http.StatusConflict,
			CodeDuplicateEmail,
			"email already exists",
		)

	// =========================
	// task
	// =========================

	case errors.Is(
		err,
		apperr.ErrTaskNotFound,
	):

		WriteError(
			w,
			http.StatusNotFound,
			CodeTaskNotFound,
			"task not found",
		)

	case errors.Is(
		err,
		apperr.ErrInvalidUserID,
	):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeInvalidUserID,
			"invalid user_id",
		)

	case errors.Is(
		err,
		apperr.ErrInvalidStatus,
	):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeValidationError,
			"invalid status",
		)

	case errors.Is(
		err,
		apperr.ErrInvalidSort,
	):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeValidationError,
			"invalid sort",
		)

	case errors.Is(
		err,
		apperr.ErrInvalidOrder,
	):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeValidationError,
			"invalid order",
		)

	case errors.Is(
		err,
		apperr.ErrInvalidTaskTitle,
	):

		WriteError(
			w,
			http.StatusBadRequest,
			CodeValidationError,
			"title is required",
		)

	case errors.Is(
		err,
		apperr.ErrForeignKeyViolation,
	):

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
