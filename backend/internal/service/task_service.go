package service

import (
	"context"

	"portfolio/backend/internal/models"
	"portfolio/backend/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepositoryInterface
}

func NewTaskService(
	r repository.TaskRepositoryInterface,
) *TaskService {

	return &TaskService{
		repo: r,
	}
}

// =========================
// status validation
// =========================

func isValidTaskStatus(status string) bool {

	switch status {

	case models.TaskStatusTODO,
		models.TaskStatusDOING,
		models.TaskStatusDONE:

		return true

	default:
		return false
	}
}

// =========================
// Create
// =========================

func (s *TaskService) Create(
	ctx context.Context,
	t *models.Task,
) (*models.Task, error) {

	if t.UserID <= 0 {
		return nil, ErrInvalidUserID
	}

	if !isValidTaskStatus(t.Status) {
		return nil, ErrInvalidStatus
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, t.ID)
}

// =========================
// Get
// =========================

func (s *TaskService) Get(
	ctx context.Context,
	id int64,
) (*models.Task, error) {

	if id <= 0 {
		return nil, ErrInvalidID
	}

	return s.repo.Get(ctx, id)
}

// =========================
// Update
// =========================

func (s *TaskService) Update(
	ctx context.Context,
	t *models.Task,
) (*models.Task, error) {

	if t.ID <= 0 {
		return nil, ErrInvalidID
	}

	if !isValidTaskStatus(t.Status) {
		return nil, ErrInvalidStatus
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, t.ID)
}

// =========================
// Delete
// =========================

func (s *TaskService) Delete(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return ErrInvalidID
	}

	return s.repo.Delete(ctx, id)
}

// =========================
// ListWithUser
// =========================

func (s *TaskService) ListWithUser(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.TaskWithUser, error) {

	return s.repo.ListWithUser(
		ctx,
		limit,
		offset,
	)
}
