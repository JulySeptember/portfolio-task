// internal/repository/task_repository.go

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"

	"portfolio/backend/internal/models"
)

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
	task *models.Task,
) error {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	query := `
INSERT INTO tasks (
	user_id,
	title,
	description,
	status,
	due_date
) VALUES (?, ?, ?, ?, ?)
`

	result, err := r.db.ExecContext(
		ctx,
		query,
		task.UserID,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
	)

	if err != nil {

		var mysqlErr *mysqlDriver.MySQLError

		if errors.As(err, &mysqlErr) {

			// ER_NO_REFERENCED_ROW_2
			if mysqlErr.Number == 1452 {
				return ErrForeignKeyViolation
			}
		}

		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	task.ID = id

	return nil
}

// =========================
// ListByUserID
// =========================

func (r *TaskRepository) ListByUserID(
	ctx context.Context,
	userID int64,
	query models.TaskListQuery,
) ([]models.Task, error) {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	sortColumn := "created_at"

	switch query.Sort {

	case "due_date":
		sortColumn = "due_date"

	case "created_at":
		sortColumn = "created_at"
	}

	order := "DESC"

	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	orderBy := "created_at DESC, id DESC"

	if sortColumn == "due_date" {

		orderBy = fmt.Sprintf(
			"due_date IS NULL, due_date %s, id DESC",
			order,
		)

	} else {

		orderBy = fmt.Sprintf(
			"%s %s, id DESC",
			sortColumn,
			order,
		)
	}

	queryStr := fmt.Sprintf(`
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
AND (? = '' OR status = ?)
ORDER BY %s
LIMIT ?
OFFSET ?
`,
		orderBy,
	)

	rows, err := r.db.QueryContext(
		ctx,
		queryStr,
		userID,
		query.Status,
		query.Status,
		query.Limit,
		query.Offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make(
		[]models.Task,
		0,
		query.Limit,
	)

	for rows.Next() {

		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(
			tasks,
			task,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// =========================
// Get
// =========================

func (r *TaskRepository) Get(
	ctx context.Context,
	taskID int64,
	userID int64,
) (*models.Task, error) {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	query := `
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

	var task models.Task

	err := r.db.QueryRowContext(
		ctx,
		query,
		taskID,
		userID,
	).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {

		if errors.Is(
			err,
			sql.ErrNoRows,
		) {

			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	return &task, nil
}

// =========================
// Update
// =========================

func (r *TaskRepository) Update(
	ctx context.Context,
	task *models.Task,
) error {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	query := `
UPDATE tasks
SET
	title = ?,
	description = ?,
	status = ?,
	due_date = ?,
	updated_at = CURRENT_TIMESTAMP
WHERE id = ?
AND user_id = ?
`

	result, err := r.db.ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.ID,
		task.UserID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

// =========================
// UpdateStatus
// =========================

func (r *TaskRepository) UpdateStatus(
	ctx context.Context,
	taskID int64,
	userID int64,
	status models.TaskStatus,
) error {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	query := `
UPDATE tasks
SET
	status = ?,
	updated_at = CURRENT_TIMESTAMP
WHERE id = ?
AND user_id = ?
`

	result, err := r.db.ExecContext(
		ctx,
		query,
		status,
		taskID,
		userID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

// =========================
// Delete
// =========================

func (r *TaskRepository) Delete(
	ctx context.Context,
	taskID int64,
	userID int64,
) error {

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	query := `
DELETE FROM tasks
WHERE id = ?
AND user_id = ?
`

	result, err := r.db.ExecContext(
		ctx,
		query,
		taskID,
		userID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}
