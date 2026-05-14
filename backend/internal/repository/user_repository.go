package repository

import (
	"context"
	"database/sql"
	"errors"

	"portfolio/backend/internal/models"
)

type UserRepositoryInterface interface {
	UpsertByAuthID(ctx context.Context, u *models.User) error
	Get(ctx context.Context, id int64) (*models.User, error)
	GetByAuthID(ctx context.Context, authUserID string) (*models.User, error)
	Delete(ctx context.Context, id int64) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(
	db *sql.DB,
) *UserRepository {

	return &UserRepository{
		db: db,
	}
}

// =========================
// UPSERT
// =========================

func (r *UserRepository) UpsertByAuthID(
	ctx context.Context,
	u *models.User,
) error {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	q := `
		INSERT INTO users (
			auth_user_id,
			email
		)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			email = VALUES(email)
	`

	_, err := r.db.ExecContext(
		ctx,
		q,
		u.AuthUserID,
		u.Email,
	)

	return err
}

// =========================
// Get
// =========================

func (r *UserRepository) Get(
	ctx context.Context,
	id int64,
) (*models.User, error) {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	q := `
		SELECT
			id,
			auth_user_id,
			email,
			created_at,
			updated_at
		FROM users
		WHERE id = ?
	`

	u := &models.User{}

	err := r.db.QueryRowContext(
		ctx,
		q,
		id,
	).Scan(
		&u.ID,
		&u.AuthUserID,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return u, nil
}

// =========================
// GetByAuthID
// =========================

func (r *UserRepository) GetByAuthID(
	ctx context.Context,
	authUserID string,
) (*models.User, error) {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	q := `
		SELECT
			id,
			auth_user_id,
			email,
			created_at,
			updated_at
		FROM users
		WHERE auth_user_id = ?
	`

	u := &models.User{}

	err := r.db.QueryRowContext(
		ctx,
		q,
		authUserID,
	).Scan(
		&u.ID,
		&u.AuthUserID,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return u, nil
}

// =========================
// Delete
// =========================

func (r *UserRepository) Delete(
	ctx context.Context,
	id int64,
) error {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM users WHERE id = ?`,
		id,
	)

	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()

	if aff == 0 {
		return ErrUserNotFound
	}

	return nil
}
