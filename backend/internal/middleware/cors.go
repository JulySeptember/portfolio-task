package middleware

import (
	"net/http"
	"os"
)

func CORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		origin := os.Getenv(
			"CORS_ALLOW_ORIGIN",
		)

		if origin != "" {

			w.Header().Set(
				"Access-Control-Allow-Origin",
				origin,
			)
		}

		w.Header().Set(
			"Vary",
			"Origin",
		)

		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Authorization, Content-Type",
		)

		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET,POST,PUT,PATCH,DELETE,OPTIONS",
		)

		if r.Method == http.MethodOptions {

			w.WriteHeader(
				http.StatusNoContent,
			)

			return
		}

		next.ServeHTTP(w, r)
	})
}
