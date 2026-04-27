package router

import (
    "net/http"
    "portfolio/backend/internal/handlers"
    "portfolio/backend/internal/models"
    "portfolio/backend/internal/service"
)

func RegisterTaskRoutes(mux *http.ServeMux, svc *service.TaskService) {
    h := handlers.NewBaseHandler[models.Task](svc)
    r := NewBaseRouter("/api/v1/tasks", h)
    r.Register(mux)
}
