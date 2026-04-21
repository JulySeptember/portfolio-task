package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "time"

    "portfolio/backend/internal/models"
    "portfolio/backend/internal/service"
)

type TaskHandler struct{ svc *service.TaskService }

func NewTaskHandler(s *service.TaskService) *TaskHandler { return &TaskHandler{svc: s} }

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.listTasks(w, r)
    case http.MethodPost:
        h.createTask(w, r)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}
func (h *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil || id <= 0 { http.Error(w, "invalid id", http.StatusBadRequest); return }
    switch r.Method {
    case http.MethodGet:
        h.getTask(w, r, id)
    case http.MethodPatch:
        h.updateTask(w, r, id)
    case http.MethodDelete:
        h.deleteTask(w, r, id)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}
func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
    var req models.Task
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "invalid body", http.StatusBadRequest); return }
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
    t, err := h.svc.CreateTask(ctx, &req)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    writeJSON(w, http.StatusCreated, t)
}
func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request, id int64) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
    t, err := h.svc.GetTask(ctx, id)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    writeJSON(w, http.StatusOK, t)
}
func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query(); limit := parseInt(q.Get("limit"), 50); offset := parseInt(q.Get("offset"), 0)
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
    tasks, err := h.svc.ListTasks(ctx, limit, offset)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    writeJSON(w, http.StatusOK, tasks)
}
func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request, id int64) {
    var req models.Task
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "invalid body", http.StatusBadRequest); return }
    req.ID = id
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
    t, err := h.svc.UpdateTask(ctx, &req)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    writeJSON(w, http.StatusOK, t)
}
func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request, id int64) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
    if err := h.svc.DeleteTask(ctx, id); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    w.WriteHeader(http.StatusNoContent)
}
