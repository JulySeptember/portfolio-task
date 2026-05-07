package repository

import (
        "context"
        "database/sql"
        "fmt"
        "strings"
        "time"

        "portfolio/backend/internal/models"
)

type TaskRepository interface {
        Create(ctx context.Context, t *models.Task) error
        Get(ctx context.Context, id int64) (*models.Task, error)
        List(ctx context.Context, status string, limit, offset int) ([]*models.Task, error)
        Update(ctx context.Context, t *models.Task) error
        Delete(ctx context.Context, id int64) error
}

type TaskRepo struct {
        DB *sql.DB
}

func NewTaskRepo(db *sql.DB) TaskRepository {
        return &TaskRepo{DB: db}
}

func (r *TaskRepo) Create(ctx context.Context, t *models.Task) error {
        query := `INSERT INTO tasks (user_id, title, description, status, due_date, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
        var due interface{}
        if t.DueDate != nil {
                due = t.DueDate.UTC()
        } else {
                due = nil
        }
        now := t.CreatedAt
        if now.IsZero() {
                now = time.Now().UTC()
        }
        res, err := r.DB.ExecContext(ctx, query, t.UserID, t.Title, t.Description, t.Status, due, now, now)
        if err != nil {
                return err
        }
        id, err := res.LastInsertId()
        if err != nil {
                return err
        }
        t.ID = id
        return nil
}

func (r *TaskRepo) Get(ctx context.Context, id int64) (*models.Task, error) {
        query := `SELECT id, user_id, title, description, status, due_date, created_at, updated_at
              FROM tasks WHERE id = ?`
        row := r.DB.QueryRowContext(ctx, query, id)
        var t models.Task
        var due sql.NullTime
        if err := row.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &due, &t.CreatedAt, &t.UpdatedAt); err != nil {
                if err == sql.ErrNoRows {
                        return nil, ErrNotFound
                }
                return nil, err
        }
        if due.Valid {
                tmp := due.Time.UTC()
                t.DueDate = &tmp
        }
        return &t, nil
}

func (r *TaskRepo) List(ctx context.Context, status string, limit, offset int) ([]*models.Task, error) {
        if limit <= 0 {
                limit = 20
        }
        if offset < 0 {
                offset = 0
        }

        var args []interface{}
        var where string
        if s := strings.TrimSpace(status); s != "" {
                where = "WHERE status = ?"
                args = append(args, s)
        }

        query := fmt.Sprintf(`SELECT id, user_id, title, description, status, due_date, created_at, updated_at
                           FROM tasks %s ORDER BY id DESC LIMIT ? OFFSET ?`, where)
        args = append(args, limit, offset)

        rows, err := r.DB.QueryContext(ctx, query, args...)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var res []*models.Task
        for rows.Next() {
                var t models.Task
                var due sql.NullTime
                if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &due, &t.CreatedAt, &t.UpdatedAt); err != nil {
                        return nil, err
                }
                if due.Valid {
                        tmp := due.Time.UTC()
                        t.DueDate = &tmp
                }
                res = append(res, &t)
        }
        if err := rows.Err(); err != nil {
                return nil, err
        }
        return res, nil
}

func (r *TaskRepo) Update(ctx context.Context, t *models.Task) error {
        query := `UPDATE tasks SET title = ?, description = ?, status = ?, due_date = ?, updated_at = ? WHERE id = ?`
        var due interface{}
        if t.DueDate != nil {
                due = t.DueDate.UTC()
        } else {
                due = nil
        }
        if t.UpdatedAt.IsZero() {
                t.UpdatedAt = time.Now().UTC()
        }
        res, err := r.DB.ExecContext(ctx, query, t.Title, t.Description, t.Status, due, t.UpdatedAt, t.ID)
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

func (r *TaskRepo) Delete(ctx context.Context, id int64) error {
        query := `DELETE FROM tasks WHERE id = ?`
        res, err := r.DB.ExecContext(ctx, query, id)
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
