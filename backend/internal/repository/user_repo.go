package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"portfolio/backend/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Repository[models.User]
}

type userRepo struct {
	*BaseRepository[models.User]
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	scanOne := func(row *sql.Row) (*models.User, error) {
		var u models.User
		if err := row.Scan(
			&u.ID,
			&u.Email,
			&u.DisplayName,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrUserNotFound
			}
			return nil, err
		}
		return &u, nil
	}

	scanMany := func(rows *sql.Rows) (*models.User, error) {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.DisplayName,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		return &u, nil
	}

	base := NewBaseRepository(
		db,
		"users",
		"id, email, display_name, created_at, updated_at",
		scanOne,
		scanMany,
	)

	return &userRepo{
		BaseRepository: base,
		db:             db,
	}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) error {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO users (email, display_name, created_at, updated_at) VALUES (?, ?, NOW(), NOW())",
		u.Email, u.DisplayName,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	u.ID = id
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

func (r *userRepo) Update(ctx context.Context, u *models.User) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET email = ?, display_name = ?, updated_at = NOW() WHERE id = ?",
		u.Email, u.DisplayName, u.ID,
	)
	return err
}
