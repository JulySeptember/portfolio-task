// internal/middleware/cors.go

package middleware

import (
	"net/http"
	"os"
	"strings"
)

// =========================
// allowed origins
// =========================

func isAllowedOrigin(
	origin string,
) bool {

	allowed := os.Getenv(
		"CORS_ALLOW_ORIGINS",
	)

	if allowed == "" {
		return false
	}

	for _, v := range strings.Split(
		allowed,
		",",
	) {

		if strings.TrimSpace(v) == origin {
			return true
		}
	}

	return false
}

// =========================
// cors middleware
// =========================

func CORS(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		origin := strings.TrimSpace(
			r.Header.Get("Origin"),
		)

		// =========================
		// exact origin match
		// =========================

		if origin != "" &&
			isAllowedOrigin(origin) {

			w.Header().Set(
				"Access-Control-Allow-Origin",
				origin,
			)

			// =========================
			// credentials
			// =========================

			w.Header().Set(
				"Access-Control-Allow-Credentials",
				"true",
			)
		}

		w.Header().Add(
			"Vary",
			"Origin",
		)

		// =========================
		// allowed headers
		// =========================

		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Authorization, Content-Type",
		)

		// =========================
		// allowed methods
		// =========================

		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET,POST,PUT,PATCH,DELETE,OPTIONS",
		)

		// =========================
		// preflight
		// =========================

		if r.Method == http.MethodOptions {

			w.WriteHeader(
				http.StatusNoContent,
			)

			return
		}

		next.ServeHTTP(
			w,
			r,
		)
	})
}
