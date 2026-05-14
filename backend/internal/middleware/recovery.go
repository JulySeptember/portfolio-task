package middleware

import (
	"log"
	"net/http"
	"portfolio/backend/internal/httpx"
)

func Recovery(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {

				log.Printf(
					"[PANIC] method=%s path=%s error=%v",
					r.Method,
					r.URL.Path,
					err,
				)

				httpx.WriteError(
					w,
					http.StatusInternalServerError,
					httpx.CodeInternalServerError,
					"internal server error",
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
