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

	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {

		id, ok := handlers.ParseID(r)

		if !ok {

			handlers.WriteError(
				w,
				http.StatusBadRequest,
				apierrors.CodeInvalidID,
				"invalid id",
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

	mux.HandleFunc("/api/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {

		id, ok := handlers.ParseID(r)

		if !ok {

			handlers.WriteError(
				w,
				http.StatusBadRequest,
				apierrors.CodeInvalidID,
				"invalid id",
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

	return mux
}
