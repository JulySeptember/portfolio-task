package service

import (
	"context"
)

type Repository[T any] interface {
	Create(ctx context.Context, t *T) error
	Get(ctx context.Context, id int64) (*T, error)
	List(ctx context.Context, limit, offset int) ([]*T, error)
	Update(ctx context.Context, t *T) error
	Delete(ctx context.Context, id int64) error
}

type BaseService[T any] struct {
	Repo Repository[T]
}

func NewBaseService[T any](repo Repository[T]) *BaseService[T] {
	return &BaseService[T]{Repo: repo}
}

func (s *BaseService[T]) Create(ctx context.Context, t *T) (*T, error) {
	if err := s.Repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *BaseService[T]) Get(ctx context.Context, id int64) (*T, error) {
	return s.Repo.Get(ctx, id)
}

func (s *BaseService[T]) List(ctx context.Context, limit, offset int) ([]*T, error) {
	return s.Repo.List(ctx, limit, offset)
}

func (s *BaseService[T]) Update(ctx context.Context, t *T) (*T, error) {
	if err := s.Repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *BaseService[T]) Delete(ctx context.Context, id int64) error {
	return s.Repo.Delete(ctx, id)
}
