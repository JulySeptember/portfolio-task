package middleware

import (
	"net/http"
	"strings"

	"portfolio/backend/internal/httpx"
)

func extractBearerToken(
	w http.ResponseWriter,
	r *http.Request,
) (string, bool) {

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

		return "", false
	}

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

		return "", false
	}

	token := strings.TrimSpace(
		parts[1],
	)

	if token == "" {

		httpx.WriteError(
			w,
			http.StatusUnauthorized,
			httpx.CodeUnauthorized,
			"missing bearer token",
		)

		return "", false
	}

	return token, true
}
