package service

import (
    "context"
    "time"

    "portfolio/backend/internal/models"
    "portfolio/backend/internal/repository"
)

type TaskService struct{ repo repository.TaskRepository }

func NewTaskService(r repository.TaskRepository) *TaskService { return &TaskService{repo: r} }

func (s *TaskService) CreateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
    done := make(chan error, 1)
    go func() { done <- s.repo.Create(t) }()
    select {
    case err := <-done:
        if err != nil { return nil, err }
        return t, nil
    case <-time.After(5 * time.Second):
        return nil, context.DeadlineExceeded
    }
}
func (s *TaskService) GetTask(ctx context.Context, id int64) (*models.Task, error) { return s.repo.Get(id) }
func (s *TaskService) ListTasks(ctx context.Context, limit, offset int) ([]*models.Task, error) { return s.repo.List(limit, offset) }
func (s *TaskService) UpdateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
    if err := s.repo.Update(t); err != nil { return nil, err }
    return t, nil
}
func (s *TaskService) DeleteTask(ctx context.Context, id int64) error { return s.repo.Delete(id) }
