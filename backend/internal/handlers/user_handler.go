package handlers

import (
	"net/http"

	"portfolio/backend/internal/apierrors"
	"portfolio/backend/internal/dto"
	"portfolio/backend/internal/models"
	"portfolio/backend/internal/repository"
	"portfolio/backend/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

// =========================
// Create
// =========================

func (h *UserHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.CreateUserRequest

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

	user := &models.User{
		Email:       req.Email,
		DisplayName: req.DisplayName,
	}

	res, err := h.svc.Create(
		r.Context(),
		user,
	)

	if err != nil {

		switch err {

		case repository.ErrDuplicateEmail:

			WriteError(
				w,
				http.StatusConflict,
				apierrors.CodeDuplicateEmail,
				"email already exists",
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

	response := dto.UserResponse{
		ID:          res.ID,
		Email:       res.Email,
		DisplayName: res.DisplayName,
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

func (h *UserHandler) Get(
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

		case repository.ErrUserNotFound:

			WriteError(
				w,
				http.StatusNotFound,
				apierrors.CodeUserNotFound,
				"user not found",
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

	response := dto.UserResponse{
		ID:          res.ID,
		Email:       res.Email,
		DisplayName: res.DisplayName,
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

func (h *UserHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {

	var req dto.UpdateUserRequest

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

	user := &models.User{
		ID:          id,
		Email:       req.Email,
		DisplayName: req.DisplayName,
	}

	res, err := h.svc.Update(
		r.Context(),
		user,
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

		case repository.ErrDuplicateEmail:

			WriteError(
				w,
				http.StatusConflict,
				apierrors.CodeDuplicateEmail,
				"email already exists",
			)

			return

		case repository.ErrUserNotFound:

			WriteError(
				w,
				http.StatusNotFound,
				apierrors.CodeUserNotFound,
				"user not found",
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

	response := dto.UserResponse{
		ID:          res.ID,
		Email:       res.Email,
		DisplayName: res.DisplayName,
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

func (h *UserHandler) Delete(
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

		case repository.ErrUserNotFound:

			WriteError(
				w,
				http.StatusNotFound,
				apierrors.CodeUserNotFound,
				"user not found",
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
