package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"portfolio/backend/internal/apperr"
	"portfolio/backend/internal/models"
)

type mockTaskRepository struct {
	createFunc func(
		ctx context.Context,
		task *models.Task,
	) (*models.Task, error)

	getByPublicIDFunc func(
		ctx context.Context,
		publicID string,
		userID int64,
	) (*models.Task, error)

	listByUserIDFunc func(
		ctx context.Context,
		userID int64,
		query models.TaskListQuery,
	) (*models.TaskListResult, error)

	updateFunc func(
		ctx context.Context,
		task *models.Task,
	) (*models.Task, error)

	updateStatusFunc func(
		ctx context.Context,
		publicID string,
		userID int64,
		status models.TaskStatus,
	) (*models.Task, error)

	deleteFunc func(
		ctx context.Context,
		publicID string,
		userID int64,
	) error
}

func (m *mockTaskRepository) Create(
	ctx context.Context,
	task *models.Task,
) (*models.Task, error) {

	if m.createFunc != nil {
		return m.createFunc(ctx, task)
	}

	return task, nil
}

func (m *mockTaskRepository) ListByUserID(
	ctx context.Context,
	userID int64,
	query models.TaskListQuery,
) (*models.TaskListResult, error) {

	if m.listByUserIDFunc != nil {
		return m.listByUserIDFunc(ctx, userID, query)
	}

	return &models.TaskListResult{}, nil
}

func (m *mockTaskRepository) GetByPublicID(
	ctx context.Context,
	publicID string,
	userID int64,
) (*models.Task, error) {

	if m.getByPublicIDFunc != nil {
		return m.getByPublicIDFunc(
			ctx,
			publicID,
			userID,
		)
	}

	return nil, nil
}

func (m *mockTaskRepository) Update(
	ctx context.Context,
	task *models.Task,
) (*models.Task, error) {

	if m.updateFunc != nil {
		return m.updateFunc(ctx, task)
	}

	return task, nil
}

func (m *mockTaskRepository) UpdateStatus(
	ctx context.Context,
	publicID string,
	userID int64,
	status models.TaskStatus,
) (*models.Task, error) {

	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(
			ctx,
			publicID,
			userID,
			status,
		)
	}

	return &models.Task{}, nil
}

func (m *mockTaskRepository) Delete(
	ctx context.Context,
	publicID string,
	userID int64,
) error {

	if m.deleteFunc != nil {
		return m.deleteFunc(
			ctx,
			publicID,
			userID,
		)
	}

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

func TestCreateTask_WhitespaceOnlyTitle(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	_, err := svc.CreateTask(
		context.Background(),
		1,
		"     ",
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

func TestCreateTask_TitleTooLong(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	title := strings.Repeat("a", 256)

	_, err := svc.CreateTask(
		context.Background(),
		1,
		title,
		"desc",
		models.TaskStatusTODO,
		nil,
	)

	if err != apperr.ErrInvalidTaskTitle {
		t.Fatalf("expected ErrInvalidTaskTitle")
	}
}

func TestCreateTask_DescriptionTooLong(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	description := strings.Repeat("a", 5001)

	_, err := svc.CreateTask(
		context.Background(),
		1,
		"title",
		description,
		models.TaskStatusTODO,
		nil,
	)

	if err != apperr.ErrInvalidDescription {
		t.Fatalf("expected ErrInvalidDescription")
	}
}

func TestCreateTask_InvalidDueDate(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	dueDate := time.Now().AddDate(11, 0, 0)

	_, err := svc.CreateTask(
		context.Background(),
		1,
		"title",
		"desc",
		models.TaskStatusTODO,
		&dueDate,
	)

	if err != apperr.ErrInvalidDueDate {
		t.Fatalf("expected ErrInvalidDueDate")
	}
}

func TestGetTaskByPublicID_Success(t *testing.T) {
	repo := &mockTaskRepository{
		getByPublicIDFunc: func(
			ctx context.Context,
			publicID string,
			userID int64,
		) (*models.Task, error) {

			return &models.Task{
				PublicID: publicID,
				UserID:   userID,
			}, nil
		},
	}

	svc := NewTaskService(repo)

	_, err := svc.GetTaskByPublicID(
		context.Background(),
		"task-1",
		1,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTaskByPublicID_InvalidID(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	_, err := svc.GetTaskByPublicID(
		context.Background(),
		"",
		1,
	)

	if err != apperr.ErrInvalidID {
		t.Fatalf("expected ErrInvalidID")
	}
}

func TestGetTaskByPublicID_InvalidUserID(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	_, err := svc.GetTaskByPublicID(
		context.Background(),
		"task-1",
		0,
	)

	if err != apperr.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID")
	}
}

func TestListTasks_Success(t *testing.T) {
	repo := &mockTaskRepository{
		listByUserIDFunc: func(
			ctx context.Context,
			userID int64,
			query models.TaskListQuery,
		) (*models.TaskListResult, error) {

			return &models.TaskListResult{
				Total: 1,
			}, nil
		},
	}

	svc := NewTaskService(repo)

	_, err := svc.ListTasks(
		context.Background(),
		1,
		models.TaskListQuery{},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListTasks_InvalidUserID(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	_, err := svc.ListTasks(
		context.Background(),
		0,
		models.TaskListQuery{},
	)

	if err != apperr.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID")
	}
}

func TestUpdateTask_Success(t *testing.T) {
	repo := &mockTaskRepository{
		updateFunc: func(
			ctx context.Context,
			task *models.Task,
		) (*models.Task, error) {

			if task.Title != "updated" {
				t.Fatalf("expected updated title")
			}

			return task, nil
		},
	}

	svc := NewTaskService(repo)

	_, err := svc.UpdateTask(
		context.Background(),
		"task-1",
		1,
		"updated",
		"desc",
		models.TaskStatusTODO,
		nil,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateTask_InvalidID(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	_, err := svc.UpdateTask(
		context.Background(),
		"",
		1,
		"title",
		"desc",
		models.TaskStatusTODO,
		nil,
	)

	if err != apperr.ErrInvalidID {
		t.Fatalf("expected ErrInvalidID")
	}
}

func TestUpdateTask_InvalidStatus(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	_, err := svc.UpdateTask(
		context.Background(),
		"task-1",
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

func TestUpdateStatus_Success(t *testing.T) {
	repo := &mockTaskRepository{}

	svc := NewTaskService(repo)

	_, err := svc.UpdateStatus(
		context.Background(),
		"task-1",
		1,
		models.TaskStatusDONE,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateStatus_InvalidStatus(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	_, err := svc.UpdateStatus(
		context.Background(),
		"task-1",
		1,
		models.TaskStatus("INVALID"),
	)

	if err != apperr.ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus")
	}
}

func TestDeleteTask_Success(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	err := svc.DeleteTask(
		context.Background(),
		"task-1",
		1,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteTask_InvalidID(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	err := svc.DeleteTask(
		context.Background(),
		"",
		1,
	)

	if err != apperr.ErrInvalidID {
		t.Fatalf("expected ErrInvalidID")
	}
}

func TestDeleteTask_InvalidUserID(t *testing.T) {
	svc := NewTaskService(&mockTaskRepository{})

	err := svc.DeleteTask(
		context.Background(),
		"task-1",
		0,
	)

	if err != apperr.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID")
	}
}
