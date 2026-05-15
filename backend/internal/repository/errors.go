package repository

import "errors"

// =========================
// base
// =========================

var ErrNotFound = errors.New("not found")

// =========================
// user
// =========================

var ErrUserNotFound = errors.New("user not found")
var ErrDuplicateEmail = errors.New("duplicate email")
var ErrDuplicateAuthUserID = errors.New("duplicate auth user id")

// =========================
// task
// =========================

var ErrTaskNotFound = errors.New("task not found")

// =========================
// db common
// =========================

var ErrForeignKeyViolation = errors.New("foreign key violation")
