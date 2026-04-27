package dto

type CreateUserRequest struct {
    Email       string `json:"email"`
    DisplayName string `json:"display_name"`
}

type UpdateUserRequest struct {
    DisplayName string `json:"display_name"`
}
