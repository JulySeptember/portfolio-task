// internal/repository/user_repository.go

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"portfolio/backend/internal/apperr"
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
// columns
// =========================

var userColumns = strings.Join([]string{
	"id",
	"auth_user_id",
	"email",
	"created_at",
	"updated_at",
}, ", ")

// =========================
// private get helper
// =========================

func (r *UserRepository) getByQuery(
	ctx context.Context,
	query string,
	args ...any,
) (*models.User, error) {

	user := &models.User{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		args...,
	).Scan(
		&user.ID,
		&user.AuthUserID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {

		if errors.Is(
			err,
			sql.ErrNoRows,
		) {

			return nil, apperr.ErrUserNotFound
		}

		return nil, parseMySQLError(err)
	}

	return user, nil
}

// =========================
// Upsert
// =========================

func (r *UserRepository) Upsert(
	ctx context.Context,
	u *models.User,
) (*models.User, error) {

	q := `
INSERT INTO users (
	auth_user_id,
	email
)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE
	email = VALUES(email),
	updated_at = CURRENT_TIMESTAMP
`

	_, err := r.db.ExecContext(
		ctx,
		q,
		u.AuthUserID,
		u.Email,
	)

	if err != nil {
		return nil, parseMySQLError(err)
	}

	return r.GetByAuthUserID(
		ctx,
		u.AuthUserID,
	)
}

// =========================
// GetByAuthUserID
// =========================

func (r *UserRepository) GetByAuthUserID(
	ctx context.Context,
	authUserID string,
) (*models.User, error) {

	q := fmt.Sprintf(`
SELECT
	%s
FROM users
WHERE auth_user_id = ?
`,
		userColumns,
	)

	return r.getByQuery(
		ctx,
		q,
		authUserID,
	)
}

// =========================
// Get
// =========================

func (r *UserRepository) Get(
	ctx context.Context,
	id int64,
) (*models.User, error) {

	q := fmt.Sprintf(`
SELECT
	%s
FROM users
WHERE id = ?
`,
		userColumns,
	)

	return r.getByQuery(
		ctx,
		q,
		id,
	)
}

// =========================
// Delete
// =========================

func (r *UserRepository) Delete(
	ctx context.Context,
	id int64,
) error {

	q := `
DELETE FROM users
WHERE id = ?
`

	res, err := r.db.ExecContext(
		ctx,
		q,
		id,
	)

	if err != nil {
		return parseMySQLError(err)
	}

	aff, err := res.RowsAffected()

	if err != nil {
		return parseMySQLError(err)
	}

	if aff == 0 {
		return apperr.ErrUserNotFound
	}

	return nil
}
