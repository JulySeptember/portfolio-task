package models

import "time"

const (
	TaskStatusTODO  = "TODO"
	TaskStatusDOING = "DOING"
	TaskStatusDONE  = "DONE"
)

// =========================
// User
// =========================

type User struct {
	ID         int64     `json:"id"`
	AuthUserID string    `json:"auth_user_id"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// =========================
// Task
// =========================

type Task struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
