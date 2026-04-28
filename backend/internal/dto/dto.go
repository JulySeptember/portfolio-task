package dto

// User DTOs

// CreateUserRequest is used when creating a new user.
type CreateUserRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

// UpdateUserRequest is used when updating an existing user.
type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
}

// Task DTOs

// CreateTaskRequest is used when creating a new task.
type CreateTaskRequest struct {
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date,omitempty"`
}

// UpdateTaskRequest is used when updating an existing task.
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date,omitempty"`
}
