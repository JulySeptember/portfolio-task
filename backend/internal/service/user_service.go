package service

import (
        "context"

        "portfolio/backend/internal/models"
        "portfolio/backend/internal/repository"
)

type UserService struct {
        *BaseService[models.User]
}

func NewUserService(repo repository.UserRepository) *UserService {
        return &UserService{
                BaseService: NewBaseService(repo),
        }
}

func (s *UserService) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
        return s.Create(ctx, u)
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
        return s.Get(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
        return s.List(ctx, limit, offset)
}

func (s *UserService) UpdateUser(ctx context.Context, u *models.User) (*models.User, error) {
        return s.Update(ctx, u)
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
        return s.Delete(ctx, id)
}
