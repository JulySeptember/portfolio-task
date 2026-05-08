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

type TaskRepositoryInterface interface {
	Create(ctx context.Context, t *models.Task) error
	Get(ctx context.Context, id int64) (*models.Task, error)
	Update(ctx context.Context, t *models.Task) error
	Delete(ctx context.Context, id int64) error

	ListWithUser(
		ctx context.Context,
		limit int,
		offset int,
	) ([]models.TaskWithUser, error)
}

// =========================
// repository
// =========================

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
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

		if errors.As(err, &mysqlErr) {

			switch mysqlErr.Number {

			case 1452:
				return ErrForeignKeyViolation
			}
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
) (*models.Task, error) {

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
	`

	t := &models.Task{}

	err := r.db.QueryRowContext(
		ctx,
		q,
		id,
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
			return nil, ErrNotFound
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

	q := `
		UPDATE tasks
		SET
			title = ?,
			description = ?,
			status = ?,
			due_date = ?
		WHERE id = ?
	`

	res, err := r.db.ExecContext(
		ctx,
		q,
		t.Title,
		t.Description,
		t.Status,
		t.DueDate,
		t.ID,
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

// =========================
// Delete
// =========================

func (r *TaskRepository) Delete(
	ctx context.Context,
	id int64,
) error {

	q := `
		DELETE FROM tasks
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

// =========================
// ListWithUser
// =========================

func (r *TaskRepository) ListWithUser(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.TaskWithUser, error) {

	q := `
		SELECT
			t.id,
			t.title,
			t.status,
			u.id,
			u.email
		FROM tasks t
		INNER JOIN users u
			ON t.user_id = u.id
		ORDER BY t.id DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(
		ctx,
		q,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.TaskWithUser

	for rows.Next() {

		var t models.TaskWithUser

		err := rows.Scan(
			&t.TaskID,
			&t.Title,
			&t.Status,
			&t.UserID,
			&t.UserEmail,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
