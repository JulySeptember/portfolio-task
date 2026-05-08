package repository

import "errors"

// =========================
// not found
// =========================

var (
	ErrUserNotFound = errors.New("user not found")
	ErrTaskNotFound = errors.New("task not found")
)

// =========================
// duplicate
// =========================

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// =========================
// foreign key
// =========================

var (
	ErrForeignKeyViolation = errors.New("foreign key violation")
)
