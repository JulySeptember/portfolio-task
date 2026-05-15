package service

import "errors"

var (

	// common
	ErrInvalidID     = errors.New("invalid id")
	ErrInvalidUserID = errors.New("invalid user id")
	ErrInvalidStatus = errors.New("invalid status")

	// user
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("duplicate email")

	// task
	ErrTaskNotFound = errors.New("task not found")

	// db
	ErrForeignKeyViolation = errors.New(
		"foreign key violation",
	)
)
