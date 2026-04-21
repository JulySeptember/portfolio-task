package service

import (
    "context"
    "time"

    "portfolio/backend/internal/models"
    "portfolio/backend/internal/repository"
)

type UserService struct{ repo repository.UserRepository }

func NewUserService(r repository.UserRepository) *UserService { return &UserService{repo: r} }

func (s *UserService) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
    done := make(chan error, 1)
    go func() { done <- s.repo.Create(u) }()
    select {
    case err := <-done:
        if err != nil { return nil, err }
        return u, nil
    case <-time.After(5 * time.Second):
        return nil, context.DeadlineExceeded
    }
}
func (s *UserService) GetUser(ctx context.Context, id int64) (*models.User, error) { return s.repo.Get(id) }
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) { return s.repo.List(limit, offset) }
func (s *UserService) UpdateUser(ctx context.Context, id int64, u *models.User) (*models.User, error) {
    u.ID = id
    if err := s.repo.Update(u); err != nil { return nil, err }
    return u, nil
}
func (s *UserService) DeleteUser(ctx context.Context, id int64) error { return s.repo.Delete(id) }
