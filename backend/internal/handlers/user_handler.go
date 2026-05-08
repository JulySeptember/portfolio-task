package handlers

import (
	"net/http"

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

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateUserRequest

	// JSON decode
	if err := DecodeJSON(w, r, &req); err != nil {
		WriteError(w, 400, err.Error())
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

	res, err := h.svc.Create(r.Context(), user)
	if err != nil {

		switch err {

		case repository.ErrDuplicateEmail:
			WriteError(w, 409, "email already exists")
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

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request, id int64) {

	res, err := h.svc.Get(r.Context(), id)
	if err != nil {

		switch err {

		case repository.ErrUserNotFound:
			WriteError(w, 404, "user not found")
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

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request, id int64) {

	var req dto.UpdateUserRequest

	// JSON decode
	if err := DecodeJSON(w, r, &req); err != nil {
		WriteError(w, 400, err.Error())
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

	res, err := h.svc.Update(r.Context(), user)
	if err != nil {

		switch err {

		case repository.ErrDuplicateEmail:
			WriteError(w, 409, "email already exists")
			return

		case repository.ErrUserNotFound:
			WriteError(w, 404, "user not found")
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

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request, id int64) {

	err := h.svc.Delete(r.Context(), id)
	if err != nil {

		switch err {

		case repository.ErrUserNotFound:
			WriteError(w, 404, "user not found")
			return

		default:
			WriteError(w, 500, "internal server error")
			return
		}
	}

	w.WriteHeader(204)
}
