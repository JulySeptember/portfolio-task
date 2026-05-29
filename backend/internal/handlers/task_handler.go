// internal/handlers/task_handler.go

package handlers

import (
	"net/http"

	"portfolio/backend/internal/dto"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/models"
	"portfolio/backend/internal/service"
)

type TaskHandler struct {
	taskSvc *service.TaskService
	userSvc *service.UserService
}

func NewTaskHandler(
	taskSvc *service.TaskService,
	userSvc *service.UserService,
) *TaskHandler {

	return &TaskHandler{
		taskSvc: taskSvc,
		userSvc: userSvc,
	}
}

// =========================
// Create
// =========================

func (h *TaskHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx, cancel := withTimeout(
		r,
		defaultHandlerTimeout,
	)
	defer cancel()

	r = r.WithContext(ctx)

	var req dto.CreateTaskRequest

	if !decodeAndValidate(
		w,
		r,
		&req,
	) {
		return
	}

	userID, ok := requireUserID(
		w,
		r,
		h.userSvc,
	)

	if !ok {
		return
	}

	dueDate, ok := parseOptionalDueDate(
		w,
		req.DueDate,
	)

	if !ok {
		return
	}

	task, err := h.taskSvc.CreateTask(
		r.Context(),
		userID,
		req.Title,
		req.Description,
		req.Status,
		dueDate,
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
		dto.ToTaskResponse(task),
	)
}

// =========================
// Get
// =========================

func (h *TaskHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx, cancel := withTimeout(
		r,
		defaultHandlerTimeout,
	)
	defer cancel()

	r = r.WithContext(ctx)

	publicID := r.PathValue("publicId")

	if publicID == "" {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"public id required",
		)

		return
	}

	userID, ok := requireUserID(
		w,
		r,
		h.userSvc,
	)

	if !ok {
		return
	}

	task, err := h.taskSvc.GetTaskByPublicID(
		r.Context(),
		publicID,
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
		dto.ToTaskResponse(task),
	)
}

// =========================
// List
// =========================

func (h *TaskHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx, cancel := withTimeout(
		r,
		listHandlerTimeout,
	)
	defer cancel()

	r = r.WithContext(ctx)

	userID, ok := requireUserID(
		w,
		r,
		h.userSvc,
	)

	if !ok {
		return
	}

	query := models.TaskListQuery{
		Status: models.TaskStatus(
			httpx.QueryString(
				r,
				"status",
				"",
			),
		),

		Sort: models.TaskSort(
			httpx.QueryString(
				r,
				"sort",
				"",
			),
		),

		Order: models.TaskOrder(
			httpx.QueryString(
				r,
				"order",
				"",
			),
		),

		Limit: httpx.QueryInt(
			r,
			"limit",
			20,
			1,
			100,
		),

		Offset: httpx.QueryInt(
			r,
			"offset",
			0,
			0,
			0,
		),
	}

	query.Normalize()

	if err := query.Validate(); err != nil {

		httpx.HandleError(
			w,
			err,
		)

		return
	}

	result, err := h.taskSvc.ListTasks(
		r.Context(),
		userID,
		query,
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
		len(result.Items),
	)

	for _, t := range result.Items {

		items = append(
			items,
			dto.ToTaskResponse(&t),
		)
	}

	resp := dto.TaskListResponse{
		Items:  items,
		Count:  result.Total,
		Limit:  query.Limit,
		Offset: query.Offset,
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		resp,
	)
}

// =========================
// Update
// =========================

func (h *TaskHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx, cancel := withTimeout(
		r,
		defaultHandlerTimeout,
	)
	defer cancel()

	r = r.WithContext(ctx)

	publicID := r.PathValue("publicId")

	if publicID == "" {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"public id required",
		)

		return
	}

	userID, ok := requireUserID(
		w,
		r,
		h.userSvc,
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

	dueDate, ok := parseOptionalDueDate(
		w,
		req.DueDate,
	)

	if !ok {
		return
	}

	task, err := h.taskSvc.UpdateTask(
		r.Context(),
		publicID,
		userID,
		req.Title,
		req.Description,
		req.Status,
		dueDate,
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
		dto.ToTaskResponse(task),
	)
}

// =========================
// Update Status
// =========================

func (h *TaskHandler) UpdateStatus(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx, cancel := withTimeout(
		r,
		defaultHandlerTimeout,
	)
	defer cancel()

	r = r.WithContext(ctx)

	publicID := r.PathValue("publicId")

	if publicID == "" {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"public id required",
		)

		return
	}

	userID, ok := requireUserID(
		w,
		r,
		h.userSvc,
	)

	if !ok {
		return
	}

	var req dto.UpdateTaskStatusRequest

	if !decodeAndValidate(
		w,
		r,
		&req,
	) {
		return
	}

	task, err := h.taskSvc.UpdateStatus(
		r.Context(),
		publicID,
		userID,
		req.Status,
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
		dto.ToTaskResponse(task),
	)
}

// =========================
// Delete
// =========================

func (h *TaskHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx, cancel := withTimeout(
		r,
		defaultHandlerTimeout,
	)
	defer cancel()

	r = r.WithContext(ctx)

	publicID := r.PathValue("publicId")

	if publicID == "" {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"public id required",
		)

		return
	}

	userID, ok := requireUserID(
		w,
		r,
		h.userSvc,
	)

	if !ok {
		return
	}

	err := h.taskSvc.DeleteTask(
		r.Context(),
		publicID,
		userID,
	)

	if err != nil {

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
