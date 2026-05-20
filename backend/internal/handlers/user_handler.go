// internal/handlers/user_handler.go

package handlers

import (
	"net/http"

	"portfolio/backend/internal/auth"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/service"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(
	userSvc *service.UserService,
) *UserHandler {

	return &UserHandler{
		userSvc: userSvc,
	}
}

// =========================
// bootstrap
// =========================

func (h *UserHandler) Bootstrap(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx, cancel := withTimeout(
		r,
		defaultHandlerTimeout,
	)
	defer cancel()

	r = r.WithContext(ctx)

	authUser, ok := auth.GetAuthUser(
		r.Context(),
	)

	if !ok {

		httpx.WriteError(
			w,
			http.StatusUnauthorized,
			httpx.CodeUnauthorized,
			"unauthorized",
		)

		return
	}

	user, err := h.userSvc.EnsureUser(
		r.Context(),
		authUser.Sub,
		authUser.Email,
	)

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusInternalServerError,
			httpx.CodeInternalServerError,
			"failed to bootstrap user",
		)

		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		user,
	)
}

// =========================
// me
// =========================

func (h *UserHandler) Me(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx, cancel := withTimeout(
		r,
		defaultHandlerTimeout,
	)
	defer cancel()

	r = r.WithContext(ctx)

	authUser, ok := auth.GetAuthUser(
		r.Context(),
	)

	if !ok {

		httpx.WriteError(
			w,
			http.StatusUnauthorized,
			httpx.CodeUnauthorized,
			"unauthorized",
		)

		return
	}

	user, err := h.userSvc.GetByAuthUserID(
		r.Context(),
		authUser.Sub,
	)

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusNotFound,
			httpx.CodeUserNotFound,
			"user not found",
		)

		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		user,
	)
}

// =========================
// delete me
// =========================

func (h *UserHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx, cancel := withTimeout(
		r,
		defaultHandlerTimeout,
	)
	defer cancel()

	r = r.WithContext(ctx)

	authUser, ok := auth.GetAuthUser(
		r.Context(),
	)

	if !ok {

		httpx.WriteError(
			w,
			http.StatusUnauthorized,
			httpx.CodeUnauthorized,
			"unauthorized",
		)

		return
	}

	user, err := h.userSvc.GetByAuthUserID(
		r.Context(),
		authUser.Sub,
	)

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusNotFound,
			httpx.CodeUserNotFound,
			"user not found",
		)

		return
	}

	err = h.userSvc.Delete(
		r.Context(),
		user.ID,
	)

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusInternalServerError,
			httpx.CodeInternalServerError,
			"failed to delete user",
		)

		return
	}

	w.WriteHeader(
		http.StatusNoContent,
	)
}
