package middleware

import (
	"net/http"
	"os"

	"portfolio/backend/internal/auth"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/service"
)

func isLocalDevAuthEnabled() bool {

	return os.Getenv("RUN_MODE") == "local" &&
		os.Getenv("AUTH_MODE") == "dev"
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
