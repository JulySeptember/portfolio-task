package service

import (
	"context"

	"portfolio/backend/internal/models"
	"portfolio/backend/internal/repository"
)

// UserService is a thin wrapper around BaseService[models.User].
// It exposes the unified CRUD API: Create, Get, List, Update, Delete.
type UserService struct {
	*BaseService[models.User]
}

// NewUserService constructs a UserService. The provided repository must
// implement repository.Repository[models.User] (repository.UserRepository typically does).
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		BaseService: NewBaseService(repo),
	}
}

// The following methods are provided to make the public API explicit.
// They simply delegate to the embedded BaseService methods.

func (s *UserService) Create(ctx context.Context, u *models.User) (*models.User, error) {
	return s.BaseService.Create(ctx, u)
}

func (s *UserService) Get(ctx context.Context, id int64) (*models.User, error) {
	return s.BaseService.Get(ctx, id)
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	return s.BaseService.List(ctx, limit, offset)
}

func (s *UserService) Update(ctx context.Context, u *models.User) (*models.User, error) {
	return s.BaseService.Update(ctx, u)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.BaseService.Delete(ctx, id)
}
