package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"portfolio/backend/internal/handlers"
	"strings"
	"time"
)

// =============================
// context key
// =============================

type contextKey string

const UserIDKey contextKey = "user_id"

// =============================
// Middleware type
// =============================

type Middleware func(http.Handler) http.Handler

// =============================
// Chain
// =============================

func Chain(h http.Handler, m ...Middleware) http.Handler {

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}

// =============================
// response writer
// =============================

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// =============================
// Logging
// =============================

func Logging(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.Printf(
			"[REQ] method=%s path=%s status=%d duration=%s remote=%s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// =============================
// Recovery
// panic recovery middleware
// =============================

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

				handlers.WriteError(
					w,
					http.StatusInternalServerError,
					"internal server error",
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// =============================
// CORS
// =============================

func CORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET,POST,PUT,PATCH,DELETE,OPTIONS",
		)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// =============================
// JWT middleware
// =============================

func JWT(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// =============================
		// local mode
		// =============================

		if os.Getenv("RUN_MODE") == "local" {

			ctx := context.WithValue(
				r.Context(),
				UserIDKey,
				int64(1),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// =============================
		// production mode
		// =============================

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {

			http.Error(
				w,
				"missing Authorization header",
				http.StatusUnauthorized,
			)

			return
		}

		parts := strings.Fields(authHeader)

		if len(parts) != 2 ||
			strings.ToLower(parts[0]) != "bearer" {

			http.Error(
				w,
				"invalid Authorization header format",
				http.StatusUnauthorized,
			)

			return
		}

		token := parts[1]

		userID := parseToken(token)
		if userID == 0 {

			http.Error(
				w,
				"invalid token",
				http.StatusUnauthorized,
			)

			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			userID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// =============================
// token parser (temporary)
// =============================

func parseToken(token string) int64 {

	// TODO:
	// replace with real JWT verification
	// ex: Cognito JWT validation

	if strings.HasPrefix(token, "user:") {

		var id int64

		_, err := fmt.Sscanf(
			token,
			"user:%d",
			&id,
		)

		if err == nil {
			return id
		}
	}

	return 0
}
