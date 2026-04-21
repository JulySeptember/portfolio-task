package repository

import (
    "database/sql"
    "errors"
    "time"

    "portfolio/backend/internal/models"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
    Create(t *models.Task) error
    Get(id int64) (*models.Task, error)
    List(limit, offset int) ([]*models.Task, error)
    Update(t *models.Task) error
    Delete(id int64) error
}

type taskRepo struct{ db *sql.DB }

func NewTaskRepository(db *sql.DB) TaskRepository { return &taskRepo{db: db} }

func (r *taskRepo) Create(t *models.Task) error {
    res, err := r.db.Exec("INSERT INTO tasks (user_id, title, description, status, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())", t.UserID, t.Title, t.Description, t.Status)
    if err != nil { return err }
    id, _ := res.LastInsertId()
    t.ID = id
    t.CreatedAt = time.Now()
    t.UpdatedAt = time.Now()
    return nil
}
func (r *taskRepo) Get(id int64) (*models.Task, error) {
    row := r.db.QueryRow("SELECT id, user_id, title, description, status, created_at, updated_at FROM tasks WHERE id = ?", id)
    var t models.Task
    if err := row.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
        if err == sql.ErrNoRows { return nil, ErrTaskNotFound }
        return nil, err
    }
    return &t, nil
}
func (r *taskRepo) List(limit, offset int) ([]*models.Task, error) {
    rows, err := r.db.Query("SELECT id, user_id, title, description, status, created_at, updated_at FROM tasks ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []*models.Task
    for rows.Next() {
        var t models.Task
        if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil { return nil, err }
        out = append(out, &t)
    }
    return out, nil
}
func (r *taskRepo) Update(t *models.Task) error {
    _, err := r.db.Exec("UPDATE tasks SET title = ?, description = ?, status = ?, updated_at = NOW() WHERE id = ?", t.Title, t.Description, t.Status, t.ID)
    return err
}
func (r *taskRepo) Delete(id int64) error {
    res, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", id)
    if err != nil { return err }
    n, _ := res.RowsAffected()
    if n == 0 { return ErrTaskNotFound }
    return nil
}
