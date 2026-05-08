package handlers

import (
	"net/http"
	"time"

	"portfolio/backend/internal/dto"
	"portfolio/backend/internal/models"
	"portfolio/backend/internal/repository"
	"portfolio/backend/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: s}
}

// =========================
// Create
// =========================

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateTaskRequest

	// JSON decode
	if err := DecodeJSON(w, r, &req); err != nil {
		WriteError(w, 400, err.Error())
		return
	}

	// validation
	if errs := ValidateStruct(req); errs != nil {
		WriteJSON(w, 400, map[string]interface{}{
			"errors": errs,
		})
		return
	}

	var dueDate *time.Time

	if req.DueDate != "" {

		t, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			WriteError(w, 400, "invalid due_date format")
			return
		}

		dueDate = &t
	}

	task := &models.Task{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     dueDate,
	}

	res, err := h.svc.Create(r.Context(), task)
	if err != nil {

		switch err {

		case repository.ErrForeignKeyViolation:
			WriteError(w, 400, "invalid user_id")
			return

		default:
			WriteError(w, 500, "internal server error")
			return
		}
	}

	WriteJSON(w, 201, res)
}

// =========================
// Get
// =========================

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request, id int64) {

	res, err := h.svc.Get(r.Context(), id)
	if err != nil {

		switch err {

		case repository.ErrTaskNotFound:
			WriteError(w, 404, "task not found")
			return

		default:
			WriteError(w, 500, "internal server error")
			return
		}
	}

	WriteJSON(w, 200, res)
}

// =========================
// Update
// =========================

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request, id int64) {

	var req dto.UpdateTaskRequest

	// JSON decode
	if err := DecodeJSON(w, r, &req); err != nil {
		WriteError(w, 400, err.Error())
		return
	}

	// validation
	if errs := ValidateStruct(req); errs != nil {
		WriteJSON(w, 400, map[string]interface{}{
			"errors": errs,
		})
		return
	}

	var dueDate *time.Time

	if req.DueDate != "" {

		t, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			WriteError(w, 400, "invalid due_date format")
			return
		}

		dueDate = &t
	}

	task := &models.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     dueDate,
	}

	res, err := h.svc.Update(r.Context(), task)
	if err != nil {

		switch err {

		case repository.ErrTaskNotFound:
			WriteError(w, 404, "task not found")
			return

		default:
			WriteError(w, 500, "internal server error")
			return
		}
	}

	WriteJSON(w, 200, res)
}

// =========================
// Delete
// =========================

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request, id int64) {

	err := h.svc.Delete(r.Context(), id)
	if err != nil {

		switch err {

		case repository.ErrTaskNotFound:
			WriteError(w, 404, "task not found")
			return

		default:
			WriteError(w, 500, "internal server error")
			return
		}
	}

	w.WriteHeader(204)
}

// =========================
// ListWithUser
// =========================

func (h *TaskHandler) ListWithUser(w http.ResponseWriter, r *http.Request) {

	limit := ParseIntOrDefault(
		r.URL.Query().Get("limit"),
		20,
	)

	offset := ParseIntOrDefault(
		r.URL.Query().Get("offset"),
		0,
	)

	// limit 上限
	if limit > 100 {
		limit = 100
	}

	// 不正値対策
	if limit < 1 {
		limit = 20
	}

	if offset < 0 {
		offset = 0
	}

	res, err := h.svc.ListWithUser(r.Context(), limit, offset)

	if err != nil {
		WriteError(w, 500, "internal server error")
		return
	}

	WriteJSON(w, 200, map[string]interface{}{
		"limit":  limit,
		"offset": offset,
		"count":  len(res),
		"items":  res,
	})
}
