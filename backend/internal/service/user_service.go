package service

import (
	"context"
	"errors"
	"strings"

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

func (s *UserService) Create(ctx context.Context, u *models.User) (*models.User, error) {

	// validation
	if strings.TrimSpace(u.Email) == "" {
		return nil, errors.New("email is required")
	}

	if !strings.Contains(u.Email, "@") {
		return nil, errors.New("invalid email")
	}

	if strings.TrimSpace(u.DisplayName) == "" {
		return nil, errors.New("display_name is required")
	}

	// create
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	// DBから完全データ取得
	return s.repo.Get(ctx, u.ID)
}

// =========================
// Get
// =========================

func (s *UserService) Get(ctx context.Context, id int64) (*models.User, error) {

	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	return s.repo.Get(ctx, id)
}

// =========================
// Update
// =========================

func (s *UserService) Update(ctx context.Context, u *models.User) (*models.User, error) {

	if u.ID <= 0 {
		return nil, errors.New("invalid id")
	}

	if strings.TrimSpace(u.Email) == "" {
		return nil, errors.New("email is required")
	}

	if !strings.Contains(u.Email, "@") {
		return nil, errors.New("invalid email")
	}

	if strings.TrimSpace(u.DisplayName) == "" {
		return nil, errors.New("display_name is required")
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, u.ID)
}

// =========================
// Delete
// =========================

func (s *UserService) Delete(ctx context.Context, id int64) error {

	if id <= 0 {
		return errors.New("invalid id")
	}

	return s.repo.Delete(ctx, id)
}
