package service

import (
	"context"

	"portfolio/backend/internal/models"
	"portfolio/backend/internal/repository"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(r repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: r}
}

// =========================
// ENSURE (middleware only)
// =========================
// CognitoユーザーをDBに同期する唯一の入口
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

	if err := s.repo.UpsertByAuthID(ctx, u); err != nil {
		return nil, err
	}

	return s.repo.GetByAuthID(ctx, authUserID)
}

// =========================
// Get (internal ID)
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
// Delete (self only)
// =========================

func (s *UserService) Delete(
	ctx context.Context,
	userID int64,
) error {

	if userID <= 0 {
		return ErrInvalidUserID
	}

	return s.repo.Delete(ctx, userID)
}
