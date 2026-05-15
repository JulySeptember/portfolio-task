// internal/service/user_service.go

package service

import (
	"context"
	"errors"

	"portfolio/backend/internal/models"
	"portfolio/backend/internal/repository"
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

	res, err := s.repo.Ensure(
		ctx,
		u,
	)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrDuplicateEmail,
		) {

			return nil, ErrDuplicateEmail
		}

		return nil, err
	}

	return res, nil
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

	res, err := s.repo.Get(
		ctx,
		id,
	)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrUserNotFound,
		) {

			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return res, nil
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

	err := s.repo.Delete(
		ctx,
		userID,
	)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrUserNotFound,
		) {

			return ErrUserNotFound
		}

		return err
	}

	return nil
}
