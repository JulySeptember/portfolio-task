package repository

import "errors"

var ErrNotFound = errors.New("not found")

var (
        ErrUserNotFound = ErrNotFound
        ErrTaskNotFound = ErrNotFound
)
