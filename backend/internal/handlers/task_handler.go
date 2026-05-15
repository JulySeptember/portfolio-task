// internal/handlers/task_handler.go

package handlers

import (
	"net/http"

	"portfolio/backend/internal/dto"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/service"
)

type TaskHandler struct {
	taskSvc *service.TaskService
}

func NewTaskHandler(
	taskSvc *service.TaskService,
) *TaskHandler {

	return &TaskHandler{
		taskSvc: taskSvc,
	}
}

// =========================
// Create
// =========================

func (h *TaskHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.CreateTaskRequest

	if !decodeAndValidate(
		w,
		r,
		&req,
	) {
		return
	}

	userID, ok := requireAuthUserID(
		w,
		r,
	)

	if !ok {
		return
	}

	task, ok := buildTaskFromRequest(
		w,
		&dto.UpdateTaskRequest{
			Title:       req.Title,
			Description: req.Description,
			Status:      req.Status,
			DueDate:     req.DueDate,
		},
	)

	if !ok {
		return
	}

	res, err := h.taskSvc.Create(
		r.Context(),
		userID,
		task,
	)

	if err != nil {

		httpx.HandleError(
			w,
			err,
		)

		return
	}

	httpx.WriteJSON(
		w,
		http.StatusCreated,
		dto.ToTaskResponse(res),
	)
}

// =========================
// Get
// =========================

func (h *TaskHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, ok := parseID(
		w,
		r,
	)

	if !ok {
		return
	}

	userID, ok := requireAuthUserID(
		w,
		r,
	)

	if !ok {
		return
	}

	res, err := h.taskSvc.Get(
		r.Context(),
		id,
		userID,
	)

	if err != nil {

		httpx.HandleError(
			w,
			err,
		)

		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		dto.ToTaskResponse(res),
	)
}

// =========================
// Update
// =========================

func (h *TaskHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, ok := parseID(
		w,
		r,
	)

	if !ok {
		return
	}

	userID, ok := requireAuthUserID(
		w,
		r,
	)

	if !ok {
		return
	}

	var req dto.UpdateTaskRequest

	if !decodeAndValidate(
		w,
		r,
		&req,
	) {
		return
	}

	task, ok := buildTaskFromRequest(
		w,
		&req,
	)

	if !ok {
		return
	}

	task.ID = id

	res, err := h.taskSvc.Update(
		r.Context(),
		task,
		userID,
	)

	if err != nil {

		httpx.HandleError(
			w,
			err,
		)

		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		dto.ToTaskResponse(res),
	)
}

// =========================
// Delete
// =========================

func (h *TaskHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, ok := parseID(
		w,
		r,
	)

	if !ok {
		return
	}

	userID, ok := requireAuthUserID(
		w,
		r,
	)

	if !ok {
		return
	}

	if err := h.taskSvc.Delete(
		r.Context(),
		id,
		userID,
	); err != nil {

		httpx.HandleError(
			w,
			err,
		)

		return
	}

	w.WriteHeader(
		http.StatusNoContent,
	)
}

// =========================
// List
// =========================

func (h *TaskHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := requireAuthUserID(
		w,
		r,
	)

	if !ok {
		return
	}

	limit := httpx.QueryInt(
		r,
		"limit",
		20,
		1,
		100,
	)

	offset := httpx.QueryInt(
		r,
		"offset",
		0,
		0,
		0,
	)

	res, err := h.taskSvc.List(
		r.Context(),
		userID,
		limit,
		offset,
	)

	if err != nil {

		httpx.HandleError(
			w,
			err,
		)

		return
	}

	items := make(
		[]dto.TaskResponse,
		0,
		len(res),
	)

	for _, t := range res {

		items = append(
			items,
			dto.ToTaskResponse(&t),
		)
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		dto.TaskListResponse{
			Count:  len(items),
			Items:  items,
			Limit:  limit,
			Offset: offset,
		},
	)
}
