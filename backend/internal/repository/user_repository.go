package repository

import (
	"context"
	"database/sql"
	"errors"

	"portfolio/backend/internal/models"
)

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
// Ensure
// =========================
// atomic upsert + fetch
// =========================

func (r *UserRepository) Ensure(
	ctx context.Context,
	u *models.User,
) (*models.User, error) {

	ctxWithTimeout, cancel := withTimeout(ctx)
	defer cancel()

	q := `
		INSERT INTO users (
			auth_user_id,
			email
		)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			email = VALUES(email),
			id = LAST_INSERT_ID(id)
	`

	res, err := r.db.ExecContext(
		ctxWithTimeout,
		q,
		u.AuthUserID,
		u.Email,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	// IMPORTANT:
	// pass original ctx
	return r.Get(
		ctx,
		id,
	)
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
