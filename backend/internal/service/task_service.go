// internal/service/task_service.go

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
// Create
// =========================

func (s *TaskService) Create(
	ctx context.Context,
	userID int64,
	t *models.Task,
) (*models.Task, error) {

	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	t.UserID = userID

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return s.repo.Get(
		ctx,
		t.ID,
		userID,
	)
}

// =========================
// Get
// =========================

func (s *TaskService) Get(
	ctx context.Context,
	id int64,
	userID int64,
) (*models.Task, error) {

	if id <= 0 {
		return nil, ErrInvalidID
	}

	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	return s.repo.Get(ctx, id, userID)
}

// =========================
// Update
// =========================

func (s *TaskService) Update(
	ctx context.Context,
	t *models.Task,
	userID int64,
) (*models.Task, error) {

	if t.ID <= 0 {
		return nil, ErrInvalidID
	}

	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	t.UserID = userID

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}

	return s.repo.Get(
		ctx,
		t.ID,
		userID,
	)
}

// =========================
// Delete
// =========================

func (s *TaskService) Delete(
	ctx context.Context,
	id int64,
	userID int64,
) error {

	if id <= 0 {
		return ErrInvalidID
	}

	if userID <= 0 {
		return ErrInvalidUserID
	}

	return s.repo.Delete(
		ctx,
		id,
		userID,
	)
}

// =========================
// List
// =========================

func (s *TaskService) List(
	ctx context.Context,
	userID int64,
	limit int,
	offset int,
) ([]models.Task, error) {

	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	return s.repo.ListByUserID(
		ctx,
		userID,
		limit,
		offset,
	)
}
