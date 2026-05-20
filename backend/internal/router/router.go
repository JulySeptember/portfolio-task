package router

import (
	"net/http"

	"portfolio/backend/internal/handlers"
)

func NewRouter(
	userHandler *handlers.UserHandler,
	taskHandler *handlers.TaskHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	// =========================
	// users
	// =========================

	mux.HandleFunc(
		"POST /api/v1/auth/bootstrap",
		userHandler.Bootstrap,
	)

	mux.HandleFunc(
		"GET /api/v1/users/me",
		userHandler.Me,
	)

	mux.HandleFunc(
		"DELETE /api/v1/users/me",
		userHandler.Delete,
	)

	// =========================
	// tasks
	// =========================

	mux.HandleFunc(
		"GET /api/v1/tasks",
		taskHandler.List,
	)

	mux.HandleFunc(
		"POST /api/v1/tasks",
		taskHandler.Create,
	)

	mux.HandleFunc(
		"GET /api/v1/tasks/{id}",
		taskHandler.Get,
	)

	mux.HandleFunc(
		"PUT /api/v1/tasks/{id}",
		taskHandler.Update,
	)

	mux.HandleFunc(
		"PATCH /api/v1/tasks/{id}/status",
		taskHandler.UpdateStatus,
	)

	mux.HandleFunc(
		"DELETE /api/v1/tasks/{id}",
		taskHandler.Delete,
	)

	return mux
}
