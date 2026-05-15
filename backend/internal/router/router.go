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

	// =========================
	// health
	// =========================

	mux.HandleFunc(
		"GET /health",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			w.WriteHeader(
				http.StatusOK,
			)

			_, _ = w.Write(
				[]byte("ok"),
			)
		},
	)

	// =========================
	// users
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
	// tasks
	// =========================

	// list
	mux.HandleFunc(
		"GET /api/v1/tasks",
		taskHandler.List,
	)

	// create
	mux.HandleFunc(
		"POST /api/v1/tasks",
		taskHandler.Create,
	)

	// get
	mux.HandleFunc(
		"GET /api/v1/tasks/{id}",
		taskHandler.Get,
	)

	// full update
	mux.HandleFunc(
		"PUT /api/v1/tasks/{id}",
		taskHandler.Update,
	)

	// partial update (status)
	mux.HandleFunc(
		"PATCH /api/v1/tasks/{id}/status",
		taskHandler.UpdateStatus,
	)

	// delete
	mux.HandleFunc(
		"DELETE /api/v1/tasks/{id}",
		taskHandler.Delete,
	)

	return mux
}
