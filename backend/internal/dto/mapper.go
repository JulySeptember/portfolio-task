package dto

import "portfolio/backend/internal/models"

func ToUserResponse(u *models.User) UserResponse {

	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func ToTaskResponse(t *models.Task) TaskResponse {

	return TaskResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func ToTaskWithUserResponse(
	t models.TaskWithUser,
) TaskWithUserResponse {

	return TaskWithUserResponse{
		TaskID:    t.TaskID,
		Title:     t.Title,
		Status:    t.Status,
		UserID:    t.UserID,
		UserEmail: t.UserEmail,
	}
}
