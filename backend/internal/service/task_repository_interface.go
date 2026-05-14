package service

import (
	"context"

	"portfolio/backend/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, t *models.Task) error

	Get(
		ctx context.Context,
		id int64,
		userID int64,
	) (*models.Task, error)

	Update(
		ctx context.Context,
		t *models.Task,
	) error

	Delete(
		ctx context.Context,
		id int64,
		userID int64,
	) error

	ListByUserID(
		ctx context.Context,
		userID int64,
		limit int,
		offset int,
	) ([]models.Task, error)
}
