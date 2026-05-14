// internal/repository/task_repository.go

package repository

import (
	"context"
	"database/sql"
	"errors"

	"portfolio/backend/internal/models"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

type TaskRepositoryInterface interface {
	Create(ctx context.Context, t *models.Task) error
	Get(ctx context.Context, id int64, userID int64) (*models.Task, error)
	Update(ctx context.Context, t *models.Task) error
	Delete(ctx context.Context, id int64, userID int64) error
	ListByUserID(
		ctx context.Context,
		userID int64,
		limit int,
		offset int,
	) ([]models.Task, error)
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(
	db *sql.DB,
) *TaskRepository {

	return &TaskRepository{
		db: db,
	}
}

// =========================
// Create
// =========================

func (r *TaskRepository) Create(
	ctx context.Context,
	t *models.Task,
) error {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	q := `
		INSERT INTO tasks (
			user_id,
			title,
			description,
			status,
			due_date
		)
		VALUES (?, ?, ?, ?, ?)
	`

	res, err := r.db.ExecContext(
		ctx,
		q,
		t.UserID,
		t.Title,
		t.Description,
		t.Status,
		t.DueDate,
	)

	if err != nil {

		var mysqlErr *mysqlDriver.MySQLError

		if errors.As(err, &mysqlErr) &&
			mysqlErr.Number == 1452 {

			return ErrForeignKeyViolation
		}

		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	t.ID = id

	return nil
}

// =========================
// Get
// =========================

func (r *TaskRepository) Get(
	ctx context.Context,
	id int64,
	userID int64,
) (*models.Task, error) {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	q := `
		SELECT
			id,
			user_id,
			title,
			description,
			status,
			due_date,
			created_at,
			updated_at
		FROM tasks
		WHERE id = ?
		AND user_id = ?
	`

	t := &models.Task{}

	err := r.db.QueryRowContext(
		ctx,
		q,
		id,
		userID,
	).Scan(
		&t.ID,
		&t.UserID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.DueDate,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	return t, nil
}

// =========================
// Update
// =========================

func (r *TaskRepository) Update(
	ctx context.Context,
	t *models.Task,
) error {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	q := `
		UPDATE tasks
		SET
			title = ?,
			description = ?,
			status = ?,
			due_date = ?
		WHERE id = ?
		AND user_id = ?
	`

	res, err := r.db.ExecContext(
		ctx,
		q,
		t.Title,
		t.Description,
		t.Status,
		t.DueDate,
		t.ID,
		t.UserID,
	)

	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()

	if aff == 0 {
		return ErrTaskNotFound
	}

	return nil
}

// =========================
// Delete
// =========================

func (r *TaskRepository) Delete(
	ctx context.Context,
	id int64,
	userID int64,
) error {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	q := `
		DELETE FROM tasks
		WHERE id = ?
		AND user_id = ?
	`

	res, err := r.db.ExecContext(
		ctx,
		q,
		id,
		userID,
	)

	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()

	if aff == 0 {
		return ErrTaskNotFound
	}

	return nil
}

// =========================
// List
// =========================

func (r *TaskRepository) ListByUserID(
	ctx context.Context,
	userID int64,
	limit int,
	offset int,
) ([]models.Task, error) {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	q := `
		SELECT
			id,
			user_id,
			title,
			description,
			status,
			due_date,
			created_at,
			updated_at
		FROM tasks
		WHERE user_id = ?
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(
		ctx,
		q,
		userID,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]models.Task, 0)

	for rows.Next() {

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

		result = append(result, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
