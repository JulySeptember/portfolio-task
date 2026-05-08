package router

import (
	"net/http"

	"portfolio/backend/internal/apierrors"
	"portfolio/backend/internal/handlers"
)

func NewRouter(
	userHandler *handlers.UserHandler,
	taskHandler *handlers.TaskHandler,
) http.Handler {

	mux := http.NewServeMux()

	methodNotAllowed := func(w http.ResponseWriter) {

		handlers.WriteError(
			w,
			http.StatusMethodNotAllowed,
			apierrors.CodeMethodNotAllowed,
			"method not allowed",
		)
	}

	// =========================
	// health check
	// =========================

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {

			methodNotAllowed(w)
			return
		}

		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte("ok"))
	})

	// =========================
	// users
	// =========================

	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			userHandler.Create(w, r)

		default:
			methodNotAllowed(w)
		}
	})

	RegisterIDRoutes(
		mux,
		"/api/v1/users/",
		"users",
		IDHandler{
			Get:    userHandler.Get,
			Update: userHandler.Update,
			Delete: userHandler.Delete,
		},
	)

	// =========================
	// tasks
	// =========================

	mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			taskHandler.ListWithUser(w, r)

		case http.MethodPost:
			taskHandler.Create(w, r)

		default:
			methodNotAllowed(w)
		}
	})

	RegisterIDRoutes(
		mux,
		"/api/v1/tasks/",
		"tasks",
		IDHandler{
			Get:    taskHandler.Get,
			Update: taskHandler.Update,
			Delete: taskHandler.Delete,
		},
	)

	return mux
}
