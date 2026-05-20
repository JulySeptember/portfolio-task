// internal/service/user_service.go

package service

import (
	"context"

	"portfolio/backend/internal/apperr"
	"portfolio/backend/internal/models"
)

type UserRepository interface {

	// =========================
	// upsert
	// =========================

	Upsert(
		ctx context.Context,
		u *models.User,
	) (*models.User, error)

	// =========================
	// lookup by auth user id
	// =========================

	GetByAuthUserID(
		ctx context.Context,
		authUserID string,
	) (*models.User, error)

	Get(
		ctx context.Context,
		id int64,
	) (*models.User, error)

	Delete(
		ctx context.Context,
		id int64,
	) error
}

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
// EnsureUser
// =========================

func (s *UserService) EnsureUser(
	ctx context.Context,
	authUserID string,
	email string,
) (*models.User, error) {

	if authUserID == "" {
		return nil, apperr.ErrInvalidUserID
	}

	return s.repo.Upsert(
		ctx,
		&models.User{
			AuthUserID: authUserID,
			Email:      email,
		},
	)
}

// =========================
// GetByAuthUserID
// =========================

func (s *UserService) GetByAuthUserID(
	ctx context.Context,
	authUserID string,
) (*models.User, error) {

	if authUserID == "" {
		return nil, apperr.ErrInvalidUserID
	}

	return s.repo.GetByAuthUserID(
		ctx,
		authUserID,
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
		return nil, apperr.ErrInvalidID
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
		return apperr.ErrInvalidUserID
	}

	return s.repo.Delete(
		ctx,
		userID,
	)
}
