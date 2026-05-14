package service

import (
	"context"

	"portfolio/backend/internal/models"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(
	r UserRepository,
) *UserService {

	return &UserService{
		repo: r,
	}
}

// =========================
// ENSURE (middleware only)
// =========================
// Cognito user sync entrypoint
// =========================

func (s *UserService) Ensure(
	ctx context.Context,
	authUserID string,
	email string,
) (*models.User, error) {

	u := &models.User{
		AuthUserID: authUserID,
		Email:      email,
	}

	return s.repo.Ensure(
		ctx,
		u,
	)
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

	return s.repo.Get(
		ctx,
		id,
	)
}

// =========================
// Delete
// =========================

func (s *UserService) Delete(
	ctx context.Context,
	userID int64,
) error {

	if userID <= 0 {
		return ErrInvalidUserID
	}

	return s.repo.Delete(
		ctx,
		userID,
	)
}
