package dto

import "time"

type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=255"`
}

type UpdateUserRequest struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=255"`
}

type UserResponse struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
