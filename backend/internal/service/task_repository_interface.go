package service

import (
	"context"

	"portfolio/backend/internal/models"
)

type TaskRepository interface {
	Create(
		ctx context.Context,
		task *models.Task,
	) error

	ListByUserID(
		ctx context.Context,
		userID int64,
		query models.TaskListQuery,
	) ([]models.Task, error)

	Get(
		ctx context.Context,
		taskID int64,
		userID int64,
	) (*models.Task, error)

	Update(
		ctx context.Context,
		task *models.Task,
	) error

	UpdateStatus(
		ctx context.Context,
		taskID int64,
		userID int64,
		status models.TaskStatus,
	) error

	Delete(
		ctx context.Context,
		taskID int64,
		userID int64,
	) error
}
