package repository

import "errors"

var ErrNotFound = errors.New("not found")

var (
	ErrUserNotFound = ErrNotFound
	ErrTaskNotFound = ErrNotFound
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var (
	ErrForeignKeyViolation = errors.New("foreign key violation")
)
