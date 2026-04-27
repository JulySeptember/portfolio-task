package dto

type CreateTaskRequest struct {
    UserID      int64  `json:"user_id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"`
}

type UpdateTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"`
}
