// internal/repository/task_repository.go

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"portfolio/backend/internal/apperr"
	"portfolio/backend/internal/models"

	"github.com/google/uuid"
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

var taskColumns = strings.Join([]string{
	"id",
	"public_id",
	"user_id",
	"title",
	"description",
	"status",
	"due_date",
	"created_at",
	"updated_at",
}, ", ")

// =========================
// get helper
// =========================

func (r *TaskRepository) get(
	ctx context.Context,
	taskID int64,
	userID int64,
) (*models.Task, error) {

	query := fmt.Sprintf(`
SELECT %s
FROM tasks
WHERE id = ?
AND user_id = ?
`,
		taskColumns,
	)

	var task models.Task

	err := r.db.QueryRowContext(
		ctx,
		query,
		taskID,
		userID,
	).Scan(
		&task.ID,
		&task.PublicID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrTaskNotFound
		}

		return nil, parseMySQLError(err)
	}

	return &task, nil
}

// =========================
// Create
// =========================

func (r *TaskRepository) Create(
	ctx context.Context,
	task *models.Task,
) (*models.Task, error) {

	publicID := uuid.New().String()

	query := `
INSERT INTO tasks (
	public_id,
	user_id,
	title,
	description,
	status,
	due_date
) VALUES (?, ?, ?, ?, ?, ?)
`

	res, err := r.db.ExecContext(
		ctx,
		query,
		publicID,
		task.UserID,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
	)

	if err != nil {
		return nil, parseMySQLError(err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, parseMySQLError(err)
	}

	return r.get(
		ctx,
		id,
		task.UserID,
	)
}

// =========================
// List
// =========================

func (r *TaskRepository) ListByUserID(
	ctx context.Context,
	userID int64,
	q models.TaskListQuery,
) (*models.TaskListResult, error) {

	sortColumn := q.Sort.Column()
	orderSQL := q.Order.SQL()

	var (
		args  []any
		where []string
	)

	where = append(
		where,
		"user_id = ?",
	)

	args = append(
		args,
		userID,
	)

	if q.Status != "" {

		where = append(
			where,
			"status = ?",
		)

		args = append(
			args,
			q.Status,
		)
	}

	whereClause := strings.Join(
		where,
		" AND ",
	)

	// =========================
	// total count
	// =========================

	countQuery := fmt.Sprintf(`
SELECT COUNT(*)
FROM tasks
WHERE %s
`,
		whereClause,
	)

	var total int64

	if err := r.db.QueryRowContext(
		ctx,
		countQuery,
		args...,
	).Scan(&total); err != nil {

		return nil, parseMySQLError(err)
	}

	// =========================
	// order by
	// =========================

	var orderBy string

	if q.Sort == models.TaskSortDueDate {

		orderBy = fmt.Sprintf(
			"due_date IS NULL ASC, due_date %s, id %s",
			orderSQL,
			orderSQL,
		)

	} else {

		orderBy = fmt.Sprintf(
			"%s %s, id %s",
			sortColumn,
			orderSQL,
			orderSQL,
		)
	}

	query := fmt.Sprintf(`
SELECT %s
FROM tasks
WHERE %s
ORDER BY %s
LIMIT ?
OFFSET ?
`,
		taskColumns,
		whereClause,
		orderBy,
	)

	listArgs := append(
		append([]any{}, args...),
		q.Limit,
		q.Offset,
	)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		listArgs...,
	)

	if err != nil {
		return nil, parseMySQLError(err)
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {

		var t models.Task

		if err := rows.Scan(
			&t.ID,
			&t.PublicID,
			&t.UserID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {

			return nil, parseMySQLError(err)
		}

		tasks = append(
			tasks,
			t,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, parseMySQLError(err)
	}

	return &models.TaskListResult{
		Items: tasks,
		Total: total,
	}, nil
}

// =========================
// Get
// =========================

func (r *TaskRepository) Get(
	ctx context.Context,
	taskID int64,
	userID int64,
) (*models.Task, error) {

	return r.get(
		ctx,
		taskID,
		userID,
	)
}

func (r *TaskRepository) GetByPublicID(ctx context.Context, publicID string, userID int64) (*models.Task, error) {
	query := fmt.Sprintf(`SELECT %s FROM tasks WHERE public_id = ? AND user_id = ?`, taskColumns)

	var task models.Task
	err := r.db.QueryRowContext(ctx, query, publicID, userID).Scan(
		&task.ID, &task.PublicID, &task.UserID, &task.Title, &task.Description,
		&task.Status, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrTaskNotFound
		}
		return nil, parseMySQLError(err)
	}
	return &task, nil
}

// =========================
// Update
// =========================

func (r *TaskRepository) Update(
	ctx context.Context,
	task *models.Task,
) (*models.Task, error) {

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

	res, err := r.db.ExecContext(
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
		return nil, parseMySQLError(err)
	}

	aff, err := res.RowsAffected()

	if err != nil {
		return nil, parseMySQLError(err)
	}

	if aff == 0 {
		return nil, apperr.ErrTaskNotFound
	}

	return r.get(
		ctx,
		task.ID,
		task.UserID,
	)
}

// =========================
// UpdateStatus
// =========================

func (r *TaskRepository) UpdateStatus(
	ctx context.Context,
	taskID int64,
	userID int64,
	status models.TaskStatus,
) (*models.Task, error) {

	query := `
UPDATE tasks
SET
	status = ?,
	updated_at = CURRENT_TIMESTAMP
WHERE id = ?
AND user_id = ?
`

	res, err := r.db.ExecContext(
		ctx,
		query,
		status,
		taskID,
		userID,
	)

	if err != nil {
		return nil, parseMySQLError(err)
	}

	aff, err := res.RowsAffected()

	if err != nil {
		return nil, parseMySQLError(err)
	}

	if aff == 0 {
		return nil, apperr.ErrTaskNotFound
	}

	return r.get(
		ctx,
		taskID,
		userID,
	)
}

// =========================
// Delete
// =========================

func (r *TaskRepository) Delete(
	ctx context.Context,
	taskID int64,
	userID int64,
) error {

	res, err := r.db.ExecContext(
		ctx,
		`
DELETE FROM tasks
WHERE id = ?
AND user_id = ?
`,
		taskID,
		userID,
	)

	if err != nil {
		return parseMySQLError(err)
	}

	aff, err := res.RowsAffected()

	if err != nil {
		return parseMySQLError(err)
	}

	if aff == 0 {
		return apperr.ErrTaskNotFound
	}

	return nil
}
