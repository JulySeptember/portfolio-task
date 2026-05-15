// internal/service/task_service.go

package service

import (
	"context"
	"errors"

	"portfolio/backend/internal/models"
	"portfolio/backend/internal/repository"
)

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(
	r TaskRepository,
) *TaskService {

	return &TaskService{
		repo: r,
	}
}

// =========================
// business validation
// =========================

func validateTaskStatus(
	status models.TaskStatus,
) error {

	switch status {

	case models.TaskStatusTODO,
		models.TaskStatusDOING,
		models.TaskStatusDONE:

		return nil

	default:
		return ErrInvalidStatus
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

	if err := validateTaskStatus(
		t.Status,
	); err != nil {

		return nil, err
	}

	t.UserID = userID

	if err := s.repo.Create(
		ctx,
		t,
	); err != nil {

		if errors.Is(
			err,
			repository.ErrForeignKeyViolation,
		) {

			return nil, ErrForeignKeyViolation
		}

		return nil, err
	}

	res, err := s.repo.Get(
		ctx,
		t.ID,
		userID,
	)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrTaskNotFound,
		) {

			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	return res, nil
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

	res, err := s.repo.Get(
		ctx,
		id,
		userID,
	)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrTaskNotFound,
		) {

			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	return res, nil
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

	if err := validateTaskStatus(
		t.Status,
	); err != nil {

		return nil, err
	}

	t.UserID = userID

	if err := s.repo.Update(
		ctx,
		t,
	); err != nil {

		if errors.Is(
			err,
			repository.ErrTaskNotFound,
		) {

			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	res, err := s.repo.Get(
		ctx,
		t.ID,
		userID,
	)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrTaskNotFound,
		) {

			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	return res, nil
}

// =========================
// UpdateStatus
// =========================

func (s *TaskService) UpdateStatus(
	ctx context.Context,
	taskID int64,
	userID int64,
	status models.TaskStatus,
) (*models.Task, error) {

	if taskID <= 0 {
		return nil, ErrInvalidID
	}

	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	if err := validateTaskStatus(
		status,
	); err != nil {

		return nil, err
	}

	if err := s.repo.UpdateStatus(
		ctx,
		taskID,
		userID,
		status,
	); err != nil {

		if errors.Is(
			err,
			repository.ErrTaskNotFound,
		) {

			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	res, err := s.repo.Get(
		ctx,
		taskID,
		userID,
	)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrTaskNotFound,
		) {

			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	return res, nil
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

	err := s.repo.Delete(
		ctx,
		id,
		userID,
	)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrTaskNotFound,
		) {

			return ErrTaskNotFound
		}

		return err
	}

	return nil
}

// =========================
// List
// =========================

func (s *TaskService) List(
	ctx context.Context,
	userID int64,
	query models.TaskListQuery,
) ([]models.Task, error) {

	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	return s.repo.ListByUserID(
		ctx,
		userID,
		query,
	)
}
