package dto

type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=255"`
}

type UpdateUserRequest struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=255"`
}

type CreateTaskRequest struct {
	UserID      int64  `json:"user_id" validate:"required"`
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Status      string `json:"status" validate:"required,oneof=TODO DOING DONE"`
	DueDate     string `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Status      string `json:"status" validate:"required,oneof=TODO DOING DONE"`
	DueDate     string `json:"due_date"`
}
