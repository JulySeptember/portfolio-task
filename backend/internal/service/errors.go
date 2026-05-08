package service

import "errors"

var (
	ErrInvalidID     = errors.New("invalid id")
	ErrInvalidUserID = errors.New("invalid user id")
)
