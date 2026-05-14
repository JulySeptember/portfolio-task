// internal/middleware/logging.go

package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"portfolio/backend/internal/auth"
	"time"
)

// =========================
// response writer
// =========================

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(
	code int,
) {

	rw.statusCode = code

	rw.ResponseWriter.WriteHeader(code)
}

// =========================
// request id
// =========================

func newRequestID() string {

	b := make([]byte, 8)

	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(b)
}

// =========================
// logging middleware
// =========================

func Logging(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		start := time.Now()

		requestID := newRequestID()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		rw.Header().Set(
			"X-Request-ID",
			requestID,
		)

		next.ServeHTTP(
			rw,
			r,
		)

		userID, _ := auth.GetUserID(
			r.Context(),
		)

		duration := time.Since(start)

		log.Printf(
			"[REQ] request_id=%s user_id=%d method=%s path=%s status=%d duration_ms=%d remote=%s user_agent=%q",
			requestID,
			userID,
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration.Milliseconds(),
			r.RemoteAddr,
			r.UserAgent(),
		)
	})
}
