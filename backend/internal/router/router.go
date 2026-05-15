// internal/router/router.go

package router

import (
	"net/http"

	"portfolio/backend/internal/handlers"
)

func NewRouter(
	userHandler *handlers.UserHandler,
	taskHandler *handlers.TaskHandler,
) http.Handler {

	mux := http.NewServeMux()

	registerUserRoutes(
		mux,
		userHandler,
	)

	registerTaskRoutes(
		mux,
		taskHandler,
	)

	return mux
}

// =========================
// user routes
// =========================

func registerUserRoutes(
	mux *http.ServeMux,
	h *handlers.UserHandler,
) {

	mux.HandleFunc(
		"GET /api/v1/users/me",
		h.Me,
	)

	mux.HandleFunc(
		"DELETE /api/v1/users/me",
		h.Delete,
	)
}

// =========================
// task routes
// =========================

func registerTaskRoutes(
	mux *http.ServeMux,
	h *handlers.TaskHandler,
) {

	// list tasks
	mux.HandleFunc(
		"GET /api/v1/tasks",
		h.List,
	)

	// create task
	mux.HandleFunc(
		"POST /api/v1/tasks",
		h.Create,
	)

	// get task
	mux.HandleFunc(
		"GET /api/v1/tasks/{id}",
		h.Get,
	)

	// update task
	mux.HandleFunc(
		"PUT /api/v1/tasks/{id}",
		h.Update,
	)

	// delete task
	mux.HandleFunc(
		"DELETE /api/v1/tasks/{id}",
		h.Delete,
	)
}
