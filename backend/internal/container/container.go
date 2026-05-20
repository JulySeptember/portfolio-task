// internal/container/container.go

package container

import (
	"database/sql"

	"portfolio/backend/internal/handlers"
	"portfolio/backend/internal/repository"
	"portfolio/backend/internal/service"
)

type Container struct {
	DB *sql.DB

	UserRepository *repository.UserRepository
	TaskRepository *repository.TaskRepository

	UserService *service.UserService
	TaskService *service.TaskService

	UserHandler *handlers.UserHandler
	TaskHandler *handlers.TaskHandler
}

func New(
	db *sql.DB,
) (*Container, error) {

	userRepo := repository.NewUserRepository(
		db,
	)

	taskRepo := repository.NewTaskRepository(
		db,
	)

	userSvc := service.NewUserService(
		userRepo,
	)

	taskSvc := service.NewTaskService(
		taskRepo,
	)

	userHandler := handlers.NewUserHandler(
		userSvc,
	)

	taskHandler := handlers.NewTaskHandler(
		taskSvc,
		userSvc,
	)

	return &Container{
		DB: db,

		UserRepository: userRepo,
		TaskRepository: taskRepo,

		UserService: userSvc,
		TaskService: taskSvc,

		UserHandler: userHandler,
		TaskHandler: taskHandler,
	}, nil
}
