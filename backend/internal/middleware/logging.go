// internal/middleware/logging.go

package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"portfolio/backend/internal/auth"
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
// structured log
// =========================

type requestLog struct {
	Timestamp  string `json:"timestamp"`
	Level      string `json:"level"`
	RequestID  string `json:"request_id"`
	UserID     int64  `json:"user_id,omitempty"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Status     int    `json:"status"`
	DurationMS int64  `json:"duration_ms"`
	RemoteAddr string `json:"remote_addr"`
	UserAgent  string `json:"user_agent"`
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

		entry := requestLog{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     "INFO",
			RequestID: requestID,
			UserID:    userID,
			Method:    r.Method,
			Path:      r.URL.Path,
			Status:    rw.statusCode,
			DurationMS: time.Since(start).
				Milliseconds(),
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
		}

		b, err := json.Marshal(entry)

		if err != nil {

			log.Printf(
				`{"level":"ERROR","message":"failed to marshal request log"}`,
			)

			return
		}

		log.Println(string(b))
	})
}
