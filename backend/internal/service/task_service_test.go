package service

import (
	"context"
	"testing"

	"portfolio/backend/internal/apperr"
	"portfolio/backend/internal/models"
)

type mockTaskRepository struct {
	createFunc func(
		ctx context.Context,
		task *models.Task,
	) (*models.Task, error)
}

func (m *mockTaskRepository) Create(
	ctx context.Context,
	task *models.Task,
) (*models.Task, error) {
	return m.createFunc(ctx, task)
}

func (m *mockTaskRepository) ListByUserID(
	ctx context.Context,
	userID int64,
	query models.TaskListQuery,
) (*models.TaskListResult, error) {
	return nil, nil
}

func (m *mockTaskRepository) GetByPublicID(
	ctx context.Context,
	publicID string,
	userID int64,
) (*models.Task, error) {
	return nil, nil
}

func (m *mockTaskRepository) Update(
	ctx context.Context,
	task *models.Task,
) (*models.Task, error) {
	return nil, nil
}

func (m *mockTaskRepository) UpdateStatus(
	ctx context.Context,
	publicID string,
	userID int64,
	status models.TaskStatus,
) (*models.Task, error) {
	return nil, nil
}

func (m *mockTaskRepository) Delete(
	ctx context.Context,
	publicID string,
	userID int64,
) error {
	return nil
}

func TestCreateTask_Success(t *testing.T) {
	repo := &mockTaskRepository{
		createFunc: func(
			ctx context.Context,
			task *models.Task,
		) (*models.Task, error) {

			if task.UserID != 1 {
				t.Fatalf("expected userID=1")
			}

			if task.Title != "test title" {
				t.Fatalf("expected trimmed title")
			}

			if task.Status != models.TaskStatusTODO {
				t.Fatalf("expected TODO status")
			}

			return task, nil
		},
	}

	svc := NewTaskService(repo)

	_, err := svc.CreateTask(
		context.Background(),
		1,
		"  test title  ",
		"desc",
		"",
		nil,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateTask_InvalidUserID(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	_, err := svc.CreateTask(
		context.Background(),
		0,
		"title",
		"desc",
		models.TaskStatusTODO,
		nil,
	)

	if err != apperr.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID")
	}
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	_, err := svc.CreateTask(
		context.Background(),
		1,
		"",
		"desc",
		models.TaskStatusTODO,
		nil,
	)

	if err != apperr.ErrInvalidTaskTitle {
		t.Fatalf("expected ErrInvalidTaskTitle")
	}
}

func TestCreateTask_InvalidStatus(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	_, err := svc.CreateTask(
		context.Background(),
		1,
		"title",
		"desc",
		models.TaskStatus("INVALID"),
		nil,
	)

	if err != apperr.ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus")
	}
}
