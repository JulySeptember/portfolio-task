package middleware

import (
	"net/http"

	"portfolio/backend/internal/models"
	"portfolio/backend/internal/service"
)

func ensureAuthenticatedUser(
	r *http.Request,
	userSvc *service.UserService,
	sub string,
	email string,
) (*models.User, error) {

	return userSvc.Ensure(
		r.Context(),
		sub,
		email,
	)
}
