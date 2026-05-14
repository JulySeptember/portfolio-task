package handlers

import (
	"net/http"

	"portfolio/backend/internal/dto"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(
	s *service.UserService,
) *UserHandler {
	return &UserHandler{
		svc: s,
	}
}

// =========================
// Me
// =========================

func (h *UserHandler) Me(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}

	res, err := h.svc.Get(r.Context(), userID)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		dto.ToUserResponse(res),
	)
}

// =========================
// Delete
// =========================

func (h *UserHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}

	if err := h.svc.Delete(
		r.Context(),
		userID,
	); err != nil {

		httpx.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
