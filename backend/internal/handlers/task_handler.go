package handlers

import (
    "portfolio/backend/internal/models"
    "portfolio/backend/internal/service"
)

type TaskHandler struct {
    h *BaseHandler[models.Task]
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
    return &TaskHandler{h: NewBaseHandler(s)}
}
