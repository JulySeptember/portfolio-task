package handlers

import (
    "portfolio/backend/internal/models"
    "portfolio/backend/internal/service"
)

type UserHandler struct {
    h *BaseHandler[models.User]
}

func NewUserHandler(s *service.UserService) *UserHandler {
    return &UserHandler{h: NewBaseHandler(s)}
}
