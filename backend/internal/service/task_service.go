package service

import (
	"context"

	"portfolio/backend/internal/models"
	"portfolio/backend/internal/repository"
)

type TaskService struct {
	*BaseService[models.Task]
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		BaseService: NewBaseService(repo),
	}
}

func (s *TaskService) CreateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
	return s.Create(ctx, t)
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*models.Task, error) {
	return s.Get(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context, limit, offset int) ([]*models.Task, error) {
	return s.List(ctx, limit, offset)
}

func (s *TaskService) UpdateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
	return s.Update(ctx, t)
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	return s.Delete(ctx, id)
}
