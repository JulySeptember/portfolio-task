// internal/service/task_service.go

package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"portfolio/backend/internal/apperr"
	"portfolio/backend/internal/models"
)

const (
	maxTaskTitleLength       = 255
	maxTaskDescriptionLength = 5000
	maxDueDateYearsAhead     = 10
)

type TaskRepository interface {
	Create(
		ctx context.Context,
		task *models.Task,
	) (*models.Task, error)

	ListByUserID(
		ctx context.Context,
		userID int64,
		query models.TaskListQuery,
	) (*models.TaskListResult, error)

	Get(
		ctx context.Context,
		taskID int64,
		userID int64,
	) (*models.Task, error)

	Update(
		ctx context.Context,
		task *models.Task,
	) (*models.Task, error)

	UpdateStatus(
		ctx context.Context,
		taskID int64,
		userID int64,
		status models.TaskStatus,
	) (*models.Task, error)

	Delete(
		ctx context.Context,
		taskID int64,
		userID int64,
	) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(
	repo TaskRepository,
) *TaskService {

	return &TaskService{
		repo: repo,
	}
}

// =========================
// validation
// =========================

func validateTaskInput(
	title string,
	description string,
	dueDate *time.Time,
) error {

	if title == "" {
		return apperr.ErrInvalidTaskTitle
	}

	if utf8.RuneCountInString(title) >
		maxTaskTitleLength {

		return apperr.ErrInvalidTaskTitle
	}

	if utf8.RuneCountInString(description) >
		maxTaskDescriptionLength {

		return apperr.ErrInvalidDescription
	}

	if dueDate != nil {

		maxDueDate := time.Now().
			AddDate(maxDueDateYearsAhead, 0, 0)

		if dueDate.After(maxDueDate) {
			return apperr.ErrInvalidDueDate
		}
	}

	return nil
}

// =========================
// Create
// =========================

func (s *TaskService) CreateTask(
	ctx context.Context,
	userID int64,
	title string,
	description string,
	status models.TaskStatus,
	dueDate *time.Time,
) (*models.Task, error) {

	if userID <= 0 {
		return nil, apperr.ErrInvalidUserID
	}

	title = strings.TrimSpace(title)

	if err := validateTaskInput(
		title,
		description,
		dueDate,
	); err != nil {

		return nil, err
	}

	if status == "" {
		status = models.TaskStatusTODO
	}

	if !status.IsValid() {
		return nil, apperr.ErrInvalidStatus
	}

	task := &models.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      status,
		DueDate:     dueDate,
	}

	return s.repo.Create(
		ctx,
		task,
	)
}

// =========================
// Get
// =========================

func (s *TaskService) GetTask(
	ctx context.Context,
	id int64,
	userID int64,
) (*models.Task, error) {

	if id <= 0 {
		return nil, apperr.ErrInvalidID
	}

	if userID <= 0 {
		return nil, apperr.ErrInvalidUserID
	}

	return s.repo.Get(
		ctx,
		id,
		userID,
	)
}

// =========================
// List
// =========================

func (s *TaskService) ListTasks(
	ctx context.Context,
	userID int64,
	query models.TaskListQuery,
) (*models.TaskListResult, error) {

	if userID <= 0 {
		return nil, apperr.ErrInvalidUserID
	}

	return s.repo.ListByUserID(
		ctx,
		userID,
		query,
	)
}

// =========================
// Update
// =========================

func (s *TaskService) UpdateTask(
	ctx context.Context,
	id int64,
	userID int64,
	title string,
	description string,
	status models.TaskStatus,
	dueDate *time.Time,
) (*models.Task, error) {

	if id <= 0 {
		return nil, apperr.ErrInvalidID
	}

	if userID <= 0 {
		return nil, apperr.ErrInvalidUserID
	}

	title = strings.TrimSpace(title)

	if err := validateTaskInput(
		title,
		description,
		dueDate,
	); err != nil {

		return nil, err
	}

	if !status.IsValid() {
		return nil, apperr.ErrInvalidStatus
	}

	task := &models.Task{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      status,
		DueDate:     dueDate,
	}

	return s.repo.Update(
		ctx,
		task,
	)
}

// =========================
// Update Status
// =========================

func (s *TaskService) UpdateStatus(
	ctx context.Context,
	id int64,
	userID int64,
	status models.TaskStatus,
) (*models.Task, error) {

	if id <= 0 {
		return nil, apperr.ErrInvalidID
	}

	if userID <= 0 {
		return nil, apperr.ErrInvalidUserID
	}

	if !status.IsValid() {
		return nil, apperr.ErrInvalidStatus
	}

	return s.repo.UpdateStatus(
		ctx,
		id,
		userID,
		status,
	)
}

// =========================
// Delete
// =========================

func (s *TaskService) DeleteTask(
	ctx context.Context,
	id int64,
	userID int64,
) error {

	if id <= 0 {
		return apperr.ErrInvalidID
	}

	if userID <= 0 {
		return apperr.ErrInvalidUserID
	}

	return s.repo.Delete(
		ctx,
		id,
		userID,
	)
}
