package dto

import (
	"time"

	"portfolio/backend/internal/models"
)

// =========================
// response
// =========================

type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =========================
// mapper
// =========================

func ToUserResponse(
	u *models.User,
) UserResponse {

	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
