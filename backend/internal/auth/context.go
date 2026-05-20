// internal/auth/context.go

package auth

import "context"

type contextKey string

const authUserKey contextKey = "auth_user"

// =========================
// auth user
// =========================

type AuthUser struct {
	Sub   string
	Email string
}

func SetAuthUser(
	ctx context.Context,
	user AuthUser,
) context.Context {

	return context.WithValue(
		ctx,
		authUserKey,
		user,
	)
}

func GetAuthUser(
	ctx context.Context,
) (AuthUser, bool) {

	user, ok := ctx.Value(
		authUserKey,
	).(AuthUser)

	return user, ok
}
