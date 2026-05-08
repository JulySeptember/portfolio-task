package dto

import "time"

type CreateTaskRequest struct {
	UserID      int64  `json:"user_id" validate:"required"`
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Status      string `json:"status" validate:"required,oneof=TODO DOING DONE"`
	DueDate     string `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Status      string `json:"status" validate:"required,oneof=TODO DOING DONE"`
	DueDate     string `json:"due_date"`
}

type TaskResponse struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskWithUserResponse struct {
	TaskID    int64  `json:"task_id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	UserID    int64  `json:"user_id"`
	UserEmail string `json:"user_email"`
}

type TaskListResponse struct {
	Count  int                    `json:"count"`
	Items  []TaskWithUserResponse `json:"items"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
}
