package service

import (
	"context"

	"portfolio/backend/internal/models"
	"portfolio/backend/internal/repository"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(
	r repository.UserRepositoryInterface,
) *UserService {

	return &UserService{
		repo: r,
	}
}

// =========================
// Create
// =========================

func (s *UserService) Create(
	ctx context.Context,
	u *models.User,
) (*models.User, error) {

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, u.ID)
}

// =========================
// Get
// =========================

func (s *UserService) Get(
	ctx context.Context,
	id int64,
) (*models.User, error) {

	if id <= 0 {
		return nil, ErrInvalidID
	}

	return s.repo.Get(ctx, id)
}

// =========================
// Update
// =========================

func (s *UserService) Update(
	ctx context.Context,
	u *models.User,
) (*models.User, error) {

	if u.ID <= 0 {
		return nil, ErrInvalidID
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, u.ID)
}

// =========================
// Delete
// =========================

func (s *UserService) Delete(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return ErrInvalidID
	}

	return s.repo.Delete(ctx, id)
}
