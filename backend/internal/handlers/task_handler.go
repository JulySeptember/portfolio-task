package handlers

import (
	"net/http"

	"portfolio/backend/internal/apierrors"
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

func (h *TaskHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.CreateTaskRequest

	// JSON decode
	if err := DecodeJSON(w, r, &req); err != nil {

		WriteError(
			w,
			http.StatusBadRequest,
			apierrors.CodeInvalidJSON,
			err.Error(),
		)

		return
	}

	// validation
	if errs := ValidateStruct(req); errs != nil {

		WriteValidationErrors(w, errs)
		return
	}

	dueDate, err := ParseOptionalTime(req.DueDate)

	if err != nil {

		WriteError(
			w,
			http.StatusBadRequest,
			apierrors.CodeInvalidDueDate,
			"invalid due_date format",
		)

		return
	}

	task := &models.Task{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     dueDate,
	}

	res, err := h.svc.Create(
		r.Context(),
		task,
	)

	if err != nil {

		switch err {

		case service.ErrInvalidUserID:

			WriteError(
				w,
				http.StatusBadRequest,
				apierrors.CodeInvalidUserID,
				"invalid user_id",
			)

			return

		case repository.ErrForeignKeyViolation:

			WriteError(
				w,
				http.StatusBadRequest,
				apierrors.CodeInvalidUserID,
				"invalid user_id",
			)

			return

		default:

			WriteError(
				w,
				http.StatusInternalServerError,
				apierrors.CodeInternalServerError,
				"internal server error",
			)

			return
		}
	}

	response := dto.TaskResponse{
		ID:          res.ID,
		UserID:      res.UserID,
		Title:       res.Title,
		Description: res.Description,
		Status:      res.Status,
		DueDate:     res.DueDate,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}

	WriteJSON(
		w,
		http.StatusCreated,
		response,
	)
}

// =========================
// Get
// =========================

func (h *TaskHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {

	res, err := h.svc.Get(
		r.Context(),
		id,
	)

	if err != nil {

		switch err {

		case service.ErrInvalidID:

			WriteError(
				w,
				http.StatusBadRequest,
				apierrors.CodeInvalidID,
				"invalid id",
			)

			return

		case repository.ErrTaskNotFound:

			WriteError(
				w,
				http.StatusNotFound,
				apierrors.CodeTaskNotFound,
				"task not found",
			)

			return

		default:

			WriteError(
				w,
				http.StatusInternalServerError,
				apierrors.CodeInternalServerError,
				"internal server error",
			)

			return
		}
	}

	response := dto.TaskResponse{
		ID:          res.ID,
		UserID:      res.UserID,
		Title:       res.Title,
		Description: res.Description,
		Status:      res.Status,
		DueDate:     res.DueDate,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}

	WriteJSON(
		w,
		http.StatusOK,
		response,
	)
}

// =========================
// Update
// =========================

func (h *TaskHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {

	var req dto.UpdateTaskRequest

	if err := DecodeJSON(w, r, &req); err != nil {

		WriteError(
			w,
			http.StatusBadRequest,
			apierrors.CodeInvalidJSON,
			err.Error(),
		)

		return
	}

	if errs := ValidateStruct(req); errs != nil {

		WriteValidationErrors(w, errs)
		return
	}

	dueDate, err := ParseOptionalTime(req.DueDate)

	if err != nil {

		WriteError(
			w,
			http.StatusBadRequest,
			apierrors.CodeInvalidDueDate,
			"invalid due_date format",
		)

		return
	}

	task := &models.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     dueDate,
	}

	res, err := h.svc.Update(
		r.Context(),
		task,
	)

	if err != nil {

		switch err {

		case service.ErrInvalidID:

			WriteError(
				w,
				http.StatusBadRequest,
				apierrors.CodeInvalidID,
				"invalid id",
			)

			return

		case repository.ErrTaskNotFound:

			WriteError(
				w,
				http.StatusNotFound,
				apierrors.CodeTaskNotFound,
				"task not found",
			)

			return

		default:

			WriteError(
				w,
				http.StatusInternalServerError,
				apierrors.CodeInternalServerError,
				"internal server error",
			)

			return
		}
	}

	response := dto.TaskResponse{
		ID:          res.ID,
		UserID:      res.UserID,
		Title:       res.Title,
		Description: res.Description,
		Status:      res.Status,
		DueDate:     res.DueDate,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}

	WriteJSON(
		w,
		http.StatusOK,
		response,
	)
}

// =========================
// Delete
// =========================

func (h *TaskHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {

	err := h.svc.Delete(
		r.Context(),
		id,
	)

	if err != nil {

		switch err {

		case service.ErrInvalidID:

			WriteError(
				w,
				http.StatusBadRequest,
				apierrors.CodeInvalidID,
				"invalid id",
			)

			return

		case repository.ErrTaskNotFound:

			WriteError(
				w,
				http.StatusNotFound,
				apierrors.CodeTaskNotFound,
				"task not found",
			)

			return

		default:

			WriteError(
				w,
				http.StatusInternalServerError,
				apierrors.CodeInternalServerError,
				"internal server error",
			)

			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// =========================
// ListWithUser
// =========================

func (h *TaskHandler) ListWithUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	limit := ParseIntOrDefault(
		r.URL.Query().Get("limit"),
		20,
	)

	offset := ParseIntOrDefault(
		r.URL.Query().Get("offset"),
		0,
	)

	if limit > 100 {
		limit = 100
	}

	if limit < 1 {
		limit = 20
	}

	if offset < 0 {
		offset = 0
	}

	res, err := h.svc.ListWithUser(
		r.Context(),
		limit,
		offset,
	)

	if err != nil {

		WriteError(
			w,
			http.StatusInternalServerError,
			apierrors.CodeInternalServerError,
			"internal server error",
		)

		return
	}

	WriteJSON(
		w,
		http.StatusOK,
		dto.TaskListResponse{
			Count:  len(res),
			Items:  res,
			Limit:  limit,
			Offset: offset,
		},
	)
}
