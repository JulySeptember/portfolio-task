package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// 共通インターフェース（User/Task で共通）
type Repository[T any] interface {
	Create(ctx context.Context, t *T) error
	Get(ctx context.Context, id int64) (*T, error)
	List(ctx context.Context, limit, offset int) ([]*T, error)
	Update(ctx context.Context, t *T) error
	Delete(ctx context.Context, id int64) error
}

// Scan 関数型
type ScanOne[T any] func(*sql.Row) (*T, error)
type ScanMany[T any] func(*sql.Rows) (*T, error)

// 共通 BaseRepository
type BaseRepository[T any] struct {
	DB       *sql.DB
	Table    string
	Columns  string
	ScanOne  ScanOne[T]
	ScanMany ScanMany[T]
}

func NewBaseRepository[T any](db *sql.DB, table string, columns string, scanOne ScanOne[T], scanMany ScanMany[T]) *BaseRepository[T] {
	return &BaseRepository[T]{
		DB:       db,
		Table:    table,
		Columns:  columns,
		ScanOne:  scanOne,
		ScanMany: scanMany,
	}
}

func (r *BaseRepository[T]) Get(ctx context.Context, id int64) (*T, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", r.Columns, r.Table)
	row := r.DB.QueryRowContext(ctx, query, id)
	return r.ScanOne(row)
}

func (r *BaseRepository[T]) List(ctx context.Context, limit, offset int) ([]*T, error) {
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY id DESC LIMIT ? OFFSET ?", r.Columns, r.Table)
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*T
	for rows.Next() {
		item, err := r.ScanMany(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.Table)
	res, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}
