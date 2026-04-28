package models

import "time"

type User struct {
        ID          int64     `json:"id"`
        Email       string    `json:"email"`
        DisplayName string    `json:"display_name"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
}

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
