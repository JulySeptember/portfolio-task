package router

import (
    "net/http"
    "portfolio/backend/internal/handlers"
)

type BaseRouter[T any] struct {
    handler *handlers.BaseHandler[T]
    prefix  string
}

func NewBaseRouter[T any](prefix string, h *handlers.BaseHandler[T]) *BaseRouter[T] {
    return &BaseRouter[T]{handler: h, prefix: prefix}
}

func (r *BaseRouter[T]) Register(mux *http.ServeMux) {
    mux.HandleFunc(r.prefix, func(w http.ResponseWriter, req *http.Request) {
        switch req.Method {
        case http.MethodGet:
            r.handler.List(w, req)
        case http.MethodPost:
            r.handler.Create(w, req)
        default:
            handlers.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
        }
    })

    mux.HandleFunc(r.prefix+"/", func(w http.ResponseWriter, req *http.Request) {
        id, err := handlers.ExtractIDFromPath(req.URL.Path, r.prefix+"/")
        if err != nil || id <= 0 {
            handlers.WriteError(w, http.StatusBadRequest, "invalid id")
            return
        }

        switch req.Method {
        case http.MethodGet:
            r.handler.Get(w, req, id)
        case http.MethodPut, http.MethodPatch:
            r.handler.HandleUpdate(w, req, id)
        case http.MethodDelete:
            r.handler.Delete(w, req, id)
        default:
            handlers.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
        }
    })
}
