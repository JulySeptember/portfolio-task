package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"portfolio/backend/internal/repository"
)

// Service はハンドラが依存する汎用サービスインターフェースです。
type Service[T any] interface {
	Create(ctx context.Context, t *T) (*T, error)
	Get(ctx context.Context, id int64) (*T, error)
	List(ctx context.Context, limit, offset int) ([]*T, error)
	Update(ctx context.Context, t *T) (*T, error)
	Delete(ctx context.Context, id int64) error
}

type DTOCreateFunc[T any] func(r *http.Request) (*T, int, error)
type DTOUpdateFunc[T any] func(r *http.Request, existing *T) (int, error)

type BaseHandler[T any] struct {
	svc          Service[T]
	decodeCreate DTOCreateFunc[T]
	mergeUpdate  DTOUpdateFunc[T]
}

func NewBaseHandler[T any](svc Service[T]) *BaseHandler[T] {
	return &BaseHandler[T]{svc: svc}
}

func NewBaseHandlerWithDTO[T any](svc Service[T], dc DTOCreateFunc[T], mu DTOUpdateFunc[T]) *BaseHandler[T] {
	return &BaseHandler[T]{
		svc:          svc,
		decodeCreate: dc,
		mergeUpdate:  mu,
	}
}

func (h *BaseHandler[T]) Create(w http.ResponseWriter, r *http.Request) {
	if h.decodeCreate == nil {
		WriteError(w, http.StatusInternalServerError, "create not supported")
		return
	}

	model, status, err := h.decodeCreate(r)
	if err != nil {
		if status == 0 {
			status = http.StatusBadRequest
		}
		WriteError(w, status, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	created, err := h.svc.Create(ctx, model)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *BaseHandler[T]) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := ParseIntOrDefault(q.Get("limit"), 20)
	offset := ParseIntOrDefault(q.Get("offset"), 0)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.svc.List(ctx, limit, offset)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, items)
}

func (h *BaseHandler[T]) Get(w http.ResponseWriter, r *http.Request, id int64) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	item, err := h.svc.Get(ctx, id)
	if err != nil {
		// 汎用ハンドラでは具体的なエンティティ固有の ErrUserNotFound / ErrTaskNotFound を参照せず、
		// 共通の ErrNotFound で判定する。
		if errors.Is(err, repository.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (h *BaseHandler[T]) HandleUpdate(w http.ResponseWriter, r *http.Request, id int64) {
	if h.mergeUpdate == nil {
		WriteError(w, http.StatusInternalServerError, "update not supported")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	existing, err := h.svc.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if status, err := h.mergeUpdate(r, existing); err != nil {
		if status == 0 {
			status = http.StatusBadRequest
		}
		WriteError(w, status, err.Error())
		return
	}

	updated, err := h.svc.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, updated)
}

func (h *BaseHandler[T]) Delete(w http.ResponseWriter, r *http.Request, id int64) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
