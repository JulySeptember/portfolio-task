package router

import (
        "net/http"
        "portfolio/backend/internal/handlers"
        "portfolio/backend/internal/service"
)

func RegisterUserRoutes(mux *http.ServeMux, svc *service.UserService) {
        h := handlers.NewUserHandler(svc)
        r := NewBaseRouter("/api/v1/users", h)
        r.Register(mux)
}
