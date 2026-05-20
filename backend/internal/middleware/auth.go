// internal/middleware/auth.go

package middleware

import (
	"net/http"

	"portfolio/backend/internal/auth"
	"portfolio/backend/internal/config"
	"portfolio/backend/internal/httpx"
)

func AuthMiddleware(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		// =========================
		// local dev bypass
		// =========================

		if isLocalDevAuthEnabled() {

			handleLocalDevAuth(
				w,
				r,
				next,
			)

			return
		}

		// =========================
		// claims from api gateway
		// =========================

		sub := r.Header.Get(
			"X-Auth-Sub",
		)

		if sub == "" {

			httpx.WriteError(
				w,
				http.StatusUnauthorized,
				httpx.CodeUnauthorized,
				"missing auth subject",
			)

			return
		}

		email := r.Header.Get(
			"X-Auth-Email",
		)

		// =========================
		// auth context
		// =========================

		ctx := auth.SetAuthUser(
			r.Context(),
			auth.AuthUser{
				Sub:   sub,
				Email: email,
			},
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

// =========================
// local dev auth
// =========================

func isLocalDevAuthEnabled() bool {

	if config.RunMode() != "local" {
		return false
	}

	if !config.DevAuthBypassEnabled() {
		return false
	}

	return true
}

func handleLocalDevAuth(
	w http.ResponseWriter,
	r *http.Request,
	next http.Handler,
) {

	ctx := auth.SetAuthUser(
		r.Context(),
		auth.AuthUser{
			Sub:   "dev-user",
			Email: "dev@example.com",
		},
	)

	next.ServeHTTP(
		w,
		r.WithContext(ctx),
	)
}
