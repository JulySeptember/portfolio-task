// internal/dto/task_dto.go

package dto

import (
	"time"

	"portfolio/backend/internal/models"
)

// =========================
// request
// =========================

type CreateTaskRequest struct {
	Title       string            `json:"title" validate:"required,min=1,max=255"`
	Description string            `json:"description" validate:"max=1000"`
	Status      models.TaskStatus `json:"status" validate:"required,oneof=TODO DOING DONE"`
	DueDate     string            `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string            `json:"title" validate:"required,min=1,max=255"`
	Description string            `json:"description" validate:"max=1000"`
	Status      models.TaskStatus `json:"status" validate:"required,oneof=TODO DOING DONE"`
	DueDate     string            `json:"due_date"`
}

// =========================
// response
// =========================

type TaskResponse struct {
	ID          int64             `json:"id"`
	UserID      int64             `json:"user_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status"`
	DueDate     *time.Time        `json:"due_date"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type TaskListResponse struct {
	Count  int            `json:"count"`
	Items  []TaskResponse `json:"items"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

// =========================
// mapper
// =========================

func ToTaskResponse(
	t *models.Task,
) TaskResponse {

	return TaskResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
