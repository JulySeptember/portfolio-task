package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"portfolio/backend/internal/models"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	Repository[models.Task]
}

type taskRepo struct {
	*BaseRepository[models.Task]
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	scanOne := func(row *sql.Row) (*models.Task, error) {
		var t models.Task
		if err := row.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrTaskNotFound
			}
			return nil, err
		}
		return &t, nil
	}

	scanMany := func(rows *sql.Rows) (*models.Task, error) {
		var t models.Task
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		return &t, nil
	}

	base := NewBaseRepository(
		db,
		"tasks",
		"id, user_id, title, description, status, due_date, created_at, updated_at",
		scanOne,
		scanMany,
	)

	return &taskRepo{
		BaseRepository: base,
		db:             db,
	}
}

func (r *taskRepo) Create(ctx context.Context, t *models.Task) error {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO tasks (user_id, title, description, status, due_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())",
		t.UserID, t.Title, t.Description, t.Status, t.DueDate,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	t.ID = id
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	return nil
}

func (r *taskRepo) Update(ctx context.Context, t *models.Task) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET title = ?, description = ?, status = ?, due_date = ?, updated_at = NOW() WHERE id = ?",
		t.Title, t.Description, t.Status, t.DueDate, t.ID,
	)
	return err
}
