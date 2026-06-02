package service

import (
	"context"
	"testing"

	"portfolio/backend/internal/apperr"
	"portfolio/backend/internal/models"
)

type mockUserRepository struct {
	upsertFunc func(
		ctx context.Context,
		u *models.User,
	) (*models.User, error)

	getByAuthUserIDFunc func(
		ctx context.Context,
		authUserID string,
	) (*models.User, error)

	getFunc func(
		ctx context.Context,
		id int64,
	) (*models.User, error)

	deleteFunc func(
		ctx context.Context,
		id int64,
	) error
}

func (m *mockUserRepository) Upsert(
	ctx context.Context,
	u *models.User,
) (*models.User, error) {

	if m.upsertFunc != nil {
		return m.upsertFunc(ctx, u)
	}

	return u, nil
}

func (m *mockUserRepository) GetByAuthUserID(
	ctx context.Context,
	authUserID string,
) (*models.User, error) {

	if m.getByAuthUserIDFunc != nil {
		return m.getByAuthUserIDFunc(
			ctx,
			authUserID,
		)
	}

	return nil, nil
}

func (m *mockUserRepository) Get(
	ctx context.Context,
	id int64,
) (*models.User, error) {

	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}

	return nil, nil
}

func (m *mockUserRepository) Delete(
	ctx context.Context,
	id int64,
) error {

	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}

	return nil
}

func TestNewUserService(t *testing.T) {
	svc := NewUserService(&mockUserRepository{})

	if svc == nil {
		t.Fatal("expected service")
	}
}

func TestEnsureUser_Success(t *testing.T) {
	repo := &mockUserRepository{
		upsertFunc: func(
			ctx context.Context,
			u *models.User,
		) (*models.User, error) {

			if u.AuthUserID != "auth-123" {
				t.Fatalf("unexpected auth user id")
			}

			if u.Email != "test@example.com" {
				t.Fatalf("unexpected email")
			}

			return u, nil
		},
	}

	svc := NewUserService(repo)

	_, err := svc.EnsureUser(
		context.Background(),
		"auth-123",
		"test@example.com",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnsureUser_InvalidUserID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{})

	_, err := svc.EnsureUser(
		context.Background(),
		"",
		"test@example.com",
	)

	if err != apperr.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID")
	}
}

func TestGetByAuthUserID_Success(t *testing.T) {
	repo := &mockUserRepository{
		getByAuthUserIDFunc: func(
			ctx context.Context,
			authUserID string,
		) (*models.User, error) {

			return &models.User{
				AuthUserID: authUserID,
			}, nil
		},
	}

	svc := NewUserService(repo)

	_, err := svc.GetByAuthUserID(
		context.Background(),
		"auth-123",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetByAuthUserID_InvalidUserID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{})

	_, err := svc.GetByAuthUserID(
		context.Background(),
		"",
	)

	if err != apperr.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID")
	}
}

func TestGet_Success(t *testing.T) {
	repo := &mockUserRepository{
		getFunc: func(
			ctx context.Context,
			id int64,
		) (*models.User, error) {

			return &models.User{
				ID: id,
			}, nil
		},
	}

	svc := NewUserService(repo)

	_, err := svc.Get(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_InvalidID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{})

	_, err := svc.Get(
		context.Background(),
		0,
	)

	if err != apperr.ErrInvalidID {
		t.Fatalf("expected ErrInvalidID")
	}
}

func TestDelete_Success(t *testing.T) {
	svc := NewUserService(&mockUserRepository{})

	err := svc.Delete(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete_InvalidUserID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{})

	err := svc.Delete(
		context.Background(),
		0,
	)

	if err != apperr.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID")
	}
}
