package service

import (
        "context"
        "time"

        "portfolio/backend/internal/models"
        "portfolio/backend/internal/repository"
)

type TaskService struct {
        repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
        return &TaskService{repo: repo}
}

func (s *TaskService) Create(ctx context.Context, t *models.Task) (*models.Task, error) {
        now := time.Now().UTC()
        t.CreatedAt = now
        t.UpdatedAt = now
        if err := s.repo.Create(ctx, t); err != nil {
                return nil, err
        }
        return t, nil
}

func (s *TaskService) Get(ctx context.Context, id int64) (*models.Task, error) {
        return s.repo.Get(ctx, id)
}

func (s *TaskService) List(ctx context.Context, status string, limit, offset int) ([]*models.Task, error) {
        return s.repo.List(ctx, status, limit, offset)
}

func (s *TaskService) Update(ctx context.Context, t *models.Task) (*models.Task, error) {
        t.UpdatedAt = time.Now().UTC()
        if err := s.repo.Update(ctx, t); err != nil {
                return nil, err
        }
        return t, nil
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
        return s.repo.Delete(ctx, id)
}
