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

	// =========================
	// USER
	// =========================

	mux.HandleFunc(
		"GET /api/v1/users/me",
		userHandler.Me,
	)

	mux.HandleFunc(
		"DELETE /api/v1/users/me",
		userHandler.Delete,
	)

	// =========================
	// TASKS
	// =========================

	// list my tasks
	mux.HandleFunc(
		"GET /api/v1/tasks",
		taskHandler.List,
	)

	// create task
	mux.HandleFunc(
		"POST /api/v1/tasks",
		taskHandler.Create,
	)

	// get task
	mux.HandleFunc(
		"GET /api/v1/tasks/{id}",
		taskHandler.Get,
	)

	// update task
	mux.HandleFunc(
		"PUT /api/v1/tasks/{id}",
		taskHandler.Update,
	)

	// delete task
	mux.HandleFunc(
		"DELETE /api/v1/tasks/{id}",
		taskHandler.Delete,
	)

	return mux
}
