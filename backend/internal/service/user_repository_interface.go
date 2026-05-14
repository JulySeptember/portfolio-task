package service

import (
	"context"

	"portfolio/backend/internal/models"
)

type UserRepository interface {

	// atomic ensure
	Ensure(
		ctx context.Context,
		u *models.User,
	) (*models.User, error)

	Get(
		ctx context.Context,
		id int64,
	) (*models.User, error)

	Delete(
		ctx context.Context,
		id int64,
	) error
}
