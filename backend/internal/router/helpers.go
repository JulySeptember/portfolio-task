package router

import (
	"net/http"

	"portfolio/backend/internal/apierrors"
	"portfolio/backend/internal/handlers"
)

type IDHandler struct {
	Get    func(http.ResponseWriter, *http.Request, int64)
	Update func(http.ResponseWriter, *http.Request, int64)
	Delete func(http.ResponseWriter, *http.Request, int64)
}

func RegisterIDRoutes(
	mux *http.ServeMux,
	path string,
	resource string,
	h IDHandler,
) {

	methodNotAllowed := func(w http.ResponseWriter) {

		handlers.WriteError(
			w,
			http.StatusMethodNotAllowed,
			apierrors.CodeMethodNotAllowed,
			"method not allowed",
		)
	}

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		id, ok := handlers.ParseID(
			r,
			resource,
		)

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
			h.Get(w, r, id)

		case http.MethodPut:
			h.Update(w, r, id)

		case http.MethodDelete:
			h.Delete(w, r, id)

		default:
			methodNotAllowed(w)
		}
	})
}
