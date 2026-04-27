package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
)

type Service[T any] interface {
    Create(ctx context.Context, t *T) (*T, error)
    Get(ctx context.Context, id int64) (*T, error)
    List(ctx context.Context, limit, offset int) ([]*T, error)
    Update(ctx context.Context, t *T) (*T, error)
    Delete(ctx context.Context, id int64) error
}

type BaseHandler[T any] struct {
    svc Service[T]
}

func NewBaseHandler[T any](svc Service[T]) *BaseHandler[T] {
    return &BaseHandler[T]{svc: svc}
}

func (h *BaseHandler[T]) Create(w http.ResponseWriter, r *http.Request) {
    var req T
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        WriteError(w, http.StatusBadRequest, "invalid body")
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    created, err := h.svc.Create(ctx, &req)
    if err != nil {
        WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }

    WriteJSON(w, http.StatusCreated, created)
}

func (h *BaseHandler[T]) Get(w http.ResponseWriter, r *http.Request, id int64) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    item, err := h.svc.Get(ctx, id)
    if err != nil {
        WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }

    WriteJSON(w, http.StatusOK, item)
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

func (h *BaseHandler[T]) HandleUpdate(w http.ResponseWriter, r *http.Request, id int64) {
    var req T
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        WriteError(w, http.StatusBadRequest, "invalid body")
        return
    }

    SetID(&req, id)

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    updated, err := h.svc.Update(ctx, &req)
    if err != nil {
        WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }

    WriteJSON(w, http.StatusOK, updated)
}

func (h *BaseHandler[T]) Delete(w http.ResponseWriter, r *http.Request, id int64) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    if err := h.svc.Delete(ctx, id); err != nil {
        WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
