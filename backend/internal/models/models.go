package models

import "time"

// =========================
// task status
// =========================

const (
	TaskStatusTODO  = "TODO"
	TaskStatusDOING = "DOING"
	TaskStatusDONE  = "DONE"
)

// =========================
// User
// =========================

type User struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

// =========================
// TaskWithUser
// =========================

type TaskWithUser struct {
	TaskID    int64
	Title     string
	Status    string
	UserID    int64
	UserEmail string
}
