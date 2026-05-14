package repository

import "errors"

var ErrNotFound = errors.New("not found")

// =========================
// user
// =========================

var ErrUserNotFound = ErrNotFound

var ErrDuplicateEmail = errors.New(
	"duplicate email",
)

var ErrDuplicateAuthUserID = errors.New(
	"duplicate auth user id",
)

// =========================
// task
// =========================

var ErrTaskNotFound = ErrNotFound

// =========================
// common db
// =========================

var ErrForeignKeyViolation = errors.New(
	"foreign key violation",
)
