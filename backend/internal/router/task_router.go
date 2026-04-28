package router

import (
        "net/http"

        "portfolio/backend/internal/handlers"
        "portfolio/backend/internal/service"
)

func RegisterTaskRoutes(mux *http.ServeMux, svc *service.TaskService) {
        handler := handlers.NewTaskHandlerWrapper(svc)

        mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
                switch r.Method {
                case http.MethodGet:
                        handler.List(w, r)
                case http.MethodPost:
                        handler.Create(w, r)
                default:
                        handlers.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
                }
        })

        mux.HandleFunc("/api/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {
                id, err := handlers.ExtractIDFromPath(r.URL.Path, "/api/v1/tasks/")
                if err != nil || id <= 0 {
                        handlers.WriteError(w, http.StatusBadRequest, "invalid id")
                        return
                }

                switch r.Method {
                case http.MethodGet:
                        handler.Get(w, r, id)
                case http.MethodPut, http.MethodPatch:
                        handler.HandleUpdate(w, r, id)
                case http.MethodDelete:
                        handler.Delete(w, r, id)
                default:
                        handlers.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
                }
        })
}
