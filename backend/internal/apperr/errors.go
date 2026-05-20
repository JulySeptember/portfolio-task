package apperr

import "errors"

var (
	// =========================
	// common
	// =========================

	ErrInvalidID     = errors.New("invalid id")
	ErrInvalidUserID = errors.New("invalid user id")

	ErrInvalidStatus = errors.New("invalid status")
	ErrInvalidSort   = errors.New("invalid sort")
	ErrInvalidOrder  = errors.New("invalid order")

	ErrValidation = errors.New("validation error")
	ErrConflict   = errors.New("resource conflict")

	ErrInvalidLimit  = errors.New("invalid limit")
	ErrInvalidOffset = errors.New("invalid offset")

	// =========================
	// user
	// =========================

	ErrUserNotFound = errors.New("user not found")

	ErrDuplicateEmail      = errors.New("duplicate email")
	ErrDuplicateAuthUserID = errors.New("duplicate auth user id")

	// =========================
	// task
	// =========================

	ErrTaskNotFound       = errors.New("task not found")
	ErrInvalidTaskTitle   = errors.New("title is required")
	ErrInvalidDescription = errors.New("invalid task description")
	ErrInvalidDueDate     = errors.New("invalid task due date")

	// =========================
	// db
	// =========================

	ErrForeignKeyViolation = errors.New("foreign key violation")
)
