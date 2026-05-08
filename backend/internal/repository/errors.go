package repository

import "errors"

var ErrNotFound = errors.New("not found")

var (
	ErrUserNotFound = ErrNotFound
	ErrTaskNotFound = ErrNotFound

	ErrDuplicateEmail = errors.New("duplicate email")

	ErrForeignKeyViolation = errors.New("foreign key violation")
)
