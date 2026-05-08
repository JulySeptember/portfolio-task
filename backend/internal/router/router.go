package router

import (
	"net/http"

	"portfolio/backend/internal/handlers"
)

// NewRouter
// API routing only
func NewRouter(
	userHandler *handlers.UserHandler,
	taskHandler *handlers.TaskHandler,
) http.Handler {

	mux := http.NewServeMux()

	// =========================
	// helper
	// =========================

	methodNotAllowed := func(w http.ResponseWriter) {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
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

	// POST /api/v1/users
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			userHandler.Create(w, r)

		default:
			methodNotAllowed(w)
		}
	})

	// GET/PATCH/DELETE /api/v1/users/{id}
	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {

		id, ok := handlers.ParseID(r)
		if !ok {

			http.Error(
				w,
				"invalid id",
				http.StatusBadRequest,
			)

			return
		}

		switch r.Method {

		case http.MethodGet:
			userHandler.Get(w, r, id)

		case http.MethodPut:
			userHandler.Update(w, r, id)

		case http.MethodDelete:
			userHandler.Delete(w, r, id)

		default:
			methodNotAllowed(w)
		}
	})

	// =========================
	// tasks
	// =========================

	// GET /api/v1/tasks
	// POST /api/v1/tasks
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

	// GET/PATCH/DELETE /api/v1/tasks/{id}
	mux.HandleFunc("/api/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {

		id, ok := handlers.ParseID(r)
		if !ok {

			http.Error(
				w,
				"invalid id",
				http.StatusBadRequest,
			)

			return
		}

		switch r.Method {

		case http.MethodGet:
			taskHandler.Get(w, r, id)

		case http.MethodPut:
			taskHandler.Update(w, r, id)

		case http.MethodDelete:
			taskHandler.Delete(w, r, id)

		default:
			methodNotAllowed(w)
		}
	})

	// =========================
	// middleware chain
	// =========================

	return mux
}
