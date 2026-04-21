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

type UserHandler struct{ svc *service.UserService }

func NewUserHandler(s *service.UserService) *UserHandler { return &UserHandler{svc: s} }

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.listUsers(w, r)
    case http.MethodPost:
        h.createUser(w, r)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}
func (h *UserHandler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil || id <= 0 { http.Error(w, "invalid id", http.StatusBadRequest); return }
    switch r.Method {
    case http.MethodGet:
        h.getUser(w, r, id)
    case http.MethodPut:
        h.updateUser(w, r, id)
    case http.MethodDelete:
        h.deleteUser(w, r, id)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
    var req models.User
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "invalid body", http.StatusBadRequest); return }
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
    u, err := h.svc.CreateUser(ctx, &req)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    writeJSON(w, http.StatusCreated, u)
}
func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request, id int64) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
    u, err := h.svc.GetUser(ctx, id)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    writeJSON(w, http.StatusOK, u)
}
func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query(); limit := parseInt(q.Get("limit"), 20); offset := parseInt(q.Get("offset"), 0)
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
    users, err := h.svc.ListUsers(ctx, limit, offset)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    writeJSON(w, http.StatusOK, users)
}
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, id int64) {
    var req models.User
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "invalid body", http.StatusBadRequest); return }
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
    req.ID = id
    u, err := h.svc.UpdateUser(ctx, id, &req)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    writeJSON(w, http.StatusOK, u)
}
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request, id int64) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
    if err := h.svc.DeleteUser(ctx, id); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    w.WriteHeader(http.StatusNoContent)
}
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
    w.Header().Set("Content-Type", "application/json"); w.WriteHeader(status); _ = json.NewEncoder(w).Encode(v)
}
func parseInt(s string, def int) int { if s == "" { return def }; if v, err := strconv.Atoi(s); err == nil { return v }; return def }
