package middleware

import (
	"net/http"

	"portfolio/backend/internal/auth"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/service"
)

func AuthMiddleware(
	userSvc *service.UserService,
) Middleware {

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			// =========================
			// local bypass
			// =========================

			if isLocalDevAuthEnabled() {

				handleLocalDevAuth(
					w,
					r,
					next,
					userSvc,
				)

				return
			}

			// =========================
			// bearer token
			// =========================

			tokenString, ok := extractBearerToken(
				w,
				r,
			)

			if !ok {
				return
			}

			// =========================
			// jwt validate
			// =========================

			token, err := auth.Validate(
				tokenString,
			)

			if err != nil {

				httpx.WriteError(
					w,
					http.StatusUnauthorized,
					httpx.CodeInvalidToken,
					"invalid token",
				)

				return
			}

			// =========================
			// ensure user
			// =========================

			user, err := ensureAuthenticatedUser(
				r,
				userSvc,
				token.Sub,
				token.Email,
			)

			if err != nil {

				httpx.WriteError(
					w,
					http.StatusInternalServerError,
					httpx.CodeInternalServerError,
					"failed to ensure user",
				)

				return
			}

			// =========================
			// inject context
			// =========================

			ctx := auth.SetUserID(
				r.Context(),
				user.ID,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
