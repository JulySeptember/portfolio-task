package repository

import (
    "database/sql"
    "errors"
    "time"

    "portfolio/backend/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
    Create(u *models.User) error
    Get(id int64) (*models.User, error)
    List(limit, offset int) ([]*models.User, error)
    Update(u *models.User) error
    Delete(id int64) error
}

type userRepo struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) UserRepository { return &userRepo{db: db} }

func (r *userRepo) Create(u *models.User) error {
    res, err := r.db.Exec("INSERT INTO users (email, display_name, created_at, updated_at) VALUES (?, ?, NOW(), NOW())", u.Email, u.DisplayName)
    if err != nil { return err }
    id, _ := res.LastInsertId()
    u.ID = id
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    return nil
}
func (r *userRepo) Get(id int64) (*models.User, error) {
    row := r.db.QueryRow("SELECT id, email, display_name, created_at, updated_at FROM users WHERE id = ?", id)
    var u models.User
    if err := row.Scan(&u.ID, &u.Email, &u.DisplayName, &u.CreatedAt, &u.UpdatedAt); err != nil {
        if err == sql.ErrNoRows { return nil, ErrUserNotFound }
        return nil, err
    }
    return &u, nil
}
func (r *userRepo) List(limit, offset int) ([]*models.User, error) {
    rows, err := r.db.Query("SELECT id, email, display_name, created_at, updated_at FROM users ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []*models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Email, &u.DisplayName, &u.CreatedAt, &u.UpdatedAt); err != nil { return nil, err }
        out = append(out, &u)
    }
    return out, nil
}
func (r *userRepo) Update(u *models.User) error {
    _, err := r.db.Exec("UPDATE users SET display_name = ?, updated_at = NOW() WHERE id = ?", u.DisplayName, u.ID)
    return err
}
func (r *userRepo) Delete(id int64) error {
    _, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
    return err
}
