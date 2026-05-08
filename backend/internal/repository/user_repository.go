package repository

import (
	"context"
	"database/sql"
	"errors"

	"portfolio/backend/internal/models"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

// =========================
// interface
// =========================

type UserRepositoryInterface interface {
	Create(ctx context.Context, u *models.User) error
	Get(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id int64) error
}

// =========================
// repository
// =========================

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// =========================
// Create
// =========================

func (r *UserRepository) Create(
	ctx context.Context,
	u *models.User,
) error {

	q := `
		INSERT INTO users (
			email,
			display_name
		)
		VALUES (?, ?)
	`

	res, err := r.db.ExecContext(
		ctx,
		q,
		u.Email,
		u.DisplayName,
	)

	if err != nil {

		var mysqlErr *mysqlDriver.MySQLError

		if errors.As(err, &mysqlErr) {

			switch mysqlErr.Number {

			case 1062:
				return ErrDuplicateEmail
			}
		}

		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

// =========================
// Get
// =========================

func (r *UserRepository) Get(
	ctx context.Context,
	id int64,
) (*models.User, error) {

	q := `
		SELECT
			id,
			email,
			display_name,
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
		&u.Email,
		&u.DisplayName,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return u, nil
}

// =========================
// Update
// =========================

func (r *UserRepository) Update(
	ctx context.Context,
	u *models.User,
) error {

	q := `
		UPDATE users
		SET
			email = ?,
			display_name = ?
		WHERE id = ?
	`

	res, err := r.db.ExecContext(
		ctx,
		q,
		u.Email,
		u.DisplayName,
		u.ID,
	)

	if err != nil {

		var mysqlErr *mysqlDriver.MySQLError

		if errors.As(err, &mysqlErr) {

			switch mysqlErr.Number {

			case 1062:
				return ErrDuplicateEmail
			}
		}

		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
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
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
}
