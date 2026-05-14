package middleware

import (
	"net/http"
	"os"
	"strings"

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

			// =====================================
			// local development bypass
			// =====================================
			// local 環境のみ認証 bypass を許可
			//
			// 本番環境で AUTH_MODE=dev が
			// 誤設定されても bypass されない
			//
			// 固定 user_id=1 を直接信用せず、
			// Ensure() を通して users table と同期
			// =====================================

			if os.Getenv("RUN_MODE") == "local" &&
				os.Getenv("AUTH_MODE") == "dev" {

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

				return
			}

			// =====================================
			// Authorization header
			// =====================================

			authHeader := strings.TrimSpace(
				r.Header.Get("Authorization"),
			)

			if authHeader == "" {

				httpx.WriteError(
					w,
					http.StatusUnauthorized,
					httpx.CodeUnauthorized,
					"missing authorization header",
				)

				return
			}

			// =====================================
			// Bearer token parse
			// =====================================

			parts := strings.SplitN(
				authHeader,
				" ",
				2,
			)

			if len(parts) != 2 ||
				!strings.EqualFold(parts[0], "Bearer") {

				httpx.WriteError(
					w,
					http.StatusUnauthorized,
					httpx.CodeUnauthorized,
					"invalid authorization header",
				)

				return
			}

			tokenString := strings.TrimSpace(
				parts[1],
			)

			// =====================================
			// JWT validation
			// =====================================

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

			// =====================================
			// Ensure user in DB
			// =====================================
			// Cognito user を users table に同期
			//
			// 初回ログイン時:
			//   users table に自動作成
			//
			// 既存ユーザー:
			//   最新情報へ同期
			// =====================================

			user, err := userSvc.Ensure(
				r.Context(),
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

			// =====================================
			// set internal user id
			// =====================================

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
