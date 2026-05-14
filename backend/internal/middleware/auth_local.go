package middleware

import (
	"net/http"
	"os"

	"portfolio/backend/internal/auth"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/service"
)

func isLocalDevAuthEnabled() bool {

	// local only
	if os.Getenv("RUN_MODE") != "local" {
		return false
	}

	// explicit opt-in
	if os.Getenv("ENABLE_DEV_AUTH_BYPASS") != "true" {
		return false
	}

	// safety: never allow in production env
	if os.Getenv("APP_ENV") == "production" {
		return false
	}

	return true
}

func handleLocalDevAuth(
	w http.ResponseWriter,
	r *http.Request,
	next http.Handler,
	userSvc *service.UserService,
) {

	user, err := userSvc.Ensure(
		r.Context(),
		"dev-user",
		"dev@example.com",
	)

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusInternalServerError,
			httpx.CodeInternalServerError,
			"failed to ensure dev user",
		)

		return
	}

	ctx := auth.SetUserID(
		r.Context(),
		user.ID,
	)

	next.ServeHTTP(
		w,
		r.WithContext(ctx),
	)
}
