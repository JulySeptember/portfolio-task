package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// Repository は汎用的な CRUD インターフェースです。
// 各具象リポジトリ（User / Task 等）はこれを満たすように実装します。
type Repository[T any] interface {
	Create(ctx context.Context, t *T) error
	Get(ctx context.Context, id int64) (*T, error)
	List(ctx context.Context, limit, offset int) ([]*T, error)
	Update(ctx context.Context, t *T) error
	Delete(ctx context.Context, id int64) error
}

// ScanOne / ScanMany の型定義
type ScanOne[T any] func(*sql.Row) (*T, error)
type ScanMany[T any] func(*sql.Rows) (*T, error)

// BaseRepository は汎用的な CRUD の共通実装を提供します。
// コンストラクタで columns / scan 関数を受け取る設計にしています。
type BaseRepository[T any] struct {
	DB       *sql.DB
	Table    string
	Columns  string
	ScanOne  ScanOne[T]
	ScanMany ScanMany[T]
}

// NewBaseRepository: コンストラクタ
func NewBaseRepository[T any](db *sql.DB, table string, columns string, scanOne ScanOne[T], scanMany ScanMany[T]) *BaseRepository[T] {
	return &BaseRepository[T]{
		DB:       db,
		Table:    table,
		Columns:  columns,
		ScanOne:  scanOne,
		ScanMany: scanMany,
	}
}

// Get: 単一取得
func (r *BaseRepository[T]) Get(ctx context.Context, id int64) (*T, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", r.Columns, r.Table)
	row := r.DB.QueryRowContext(ctx, query, id)
	if r.ScanOne == nil {
		return nil, fmt.Errorf("ScanOne not implemented for table %s", r.Table)
	}
	return r.ScanOne(row)
}

// List: 一覧取得
func (r *BaseRepository[T]) List(ctx context.Context, limit, offset int) ([]*T, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY id DESC LIMIT ? OFFSET ?", r.Columns, r.Table)
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if r.ScanMany == nil {
		return nil, fmt.Errorf("ScanMany not implemented for table %s", r.Table)
	}

	var out []*T
	for rows.Next() {
		item, err := r.ScanMany(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete: 削除
func (r *BaseRepository[T]) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.Table)
	res, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
