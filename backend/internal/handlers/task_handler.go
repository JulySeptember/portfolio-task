package handlers

import (
        "context"
        "errors"
        "net/http"
        "strings"
        "time"

        "portfolio/backend/internal/dto"
        "portfolio/backend/internal/models"
        "portfolio/backend/internal/repository"
        "portfolio/backend/internal/service"
)

type TaskHandlerWrapper struct {
        svc *service.TaskService
}

func NewTaskHandlerWrapper(svc *service.TaskService) *TaskHandlerWrapper {
        return &TaskHandlerWrapper{
                svc: svc,
        }
}

func (h *TaskHandlerWrapper) List(w http.ResponseWriter, r *http.Request) {
        q := r.URL.Query()
        status := strings.ToUpper(strings.TrimSpace(q.Get("status")))
        limit := ParseIntOrDefault(q.Get("limit"), 20)
        offset := ParseIntOrDefault(q.Get("offset"), 0)

        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()

        items, err := h.svc.List(ctx, status, limit, offset)
        if err != nil {
                if errors.Is(err, repository.ErrNotFound) {
                        WriteError(w, http.StatusNotFound, "not found")
                        return
                }
                WriteError(w, http.StatusInternalServerError, err.Error())
                return
        }
        WriteJSON(w, http.StatusOK, items)
}

func (h *TaskHandlerWrapper) Create(w http.ResponseWriter, r *http.Request) {
        model, status, err := decodeCreateTaskLocal(r)
        if err != nil {
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

func (h *TaskHandlerWrapper) Get(w http.ResponseWriter, r *http.Request, id int64) {
        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()

        item, err := h.svc.Get(ctx, id)
        if err != nil {
                if errors.Is(err, repository.ErrNotFound) {
                        WriteError(w, http.StatusNotFound, "task not found")
                        return
                }
                WriteError(w, http.StatusInternalServerError, err.Error())
                return
        }
        WriteJSON(w, http.StatusOK, item)
}

func (h *TaskHandlerWrapper) HandleUpdate(w http.ResponseWriter, r *http.Request, id int64) {
        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()

        existing, err := h.svc.Get(ctx, id)
        if err != nil {
                if errors.Is(err, repository.ErrNotFound) {
                        WriteError(w, http.StatusNotFound, "task not found")
                        return
                }
                WriteError(w, http.StatusInternalServerError, err.Error())
                return
        }

        if status, err := mergeUpdateTaskLocal(r, existing); err != nil {
                WriteError(w, status, err.Error())
                return
        }

        updated, err := h.svc.Update(ctx, existing)
        if err != nil {
                if errors.Is(err, repository.ErrNotFound) {
                        WriteError(w, http.StatusNotFound, "task not found")
                        return
                }
                WriteError(w, http.StatusInternalServerError, err.Error())
                return
        }
        WriteJSON(w, http.StatusOK, updated)
}

func (h *TaskHandlerWrapper) Delete(w http.ResponseWriter, r *http.Request, id int64) {
        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()

        if err := h.svc.Delete(ctx, id); err != nil {
                if errors.Is(err, repository.ErrNotFound) {
                        WriteError(w, http.StatusNotFound, "task not found")
                        return
                }
                WriteError(w, http.StatusInternalServerError, err.Error())
                return
        }
        w.WriteHeader(http.StatusNoContent)
}

func decodeCreateTaskLocal(r *http.Request) (*models.Task, int, error) {
        var req dto.CreateTaskRequest
        if err := DecodeJSON(r, &req); err != nil {
                return nil, http.StatusBadRequest, errors.New("invalid body")
        }
        if req.UserID == 0 {
                return nil, http.StatusBadRequest, errors.New("user_id is required")
        }
        if req.Title == "" {
                return nil, http.StatusBadRequest, errors.New("title is required")
        }

        t := &models.Task{
                UserID:      req.UserID,
                Title:       req.Title,
                Description: req.Description,
                Status:      req.Status,
        }

        if req.DueDate != "" {
                tt, err := time.Parse(time.RFC3339, req.DueDate)
                if err != nil {
                        return nil, http.StatusBadRequest, errors.New("invalid due_date format; use RFC3339")
                }
                t.DueDate = &tt
        }
        return t, 0, nil
}

func mergeUpdateTaskLocal(r *http.Request, existing *models.Task) (int, error) {
        var req dto.UpdateTaskRequest
        if err := DecodeJSON(r, &req); err != nil {
                return http.StatusBadRequest, errors.New("invalid body")
        }
        if req.Title != "" {
                existing.Title = req.Title
        }
        if req.Description != "" {
                existing.Description = req.Description
        }
        if req.Status != "" {
                existing.Status = req.Status
        }
        if req.DueDate != "" {
                tt, err := time.Parse(time.RFC3339, req.DueDate)
                if err != nil {
                        return http.StatusBadRequest, errors.New("invalid due_date format; use RFC3339")
                }
                existing.DueDate = &tt
        }
        return 0, nil
}
