// internal/middleware/logging.go

package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"portfolio/backend/internal/auth"
)

// =========================
// response writer
// =========================

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

func (rw *responseWriter) WriteHeader(
	code int,
) {

	rw.statusCode = code

	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(
	b []byte,
) (int, error) {

	// implicit 200 support

	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	n, err := rw.ResponseWriter.Write(b)

	rw.bytesWritten += int64(n)

	return n, err
}

// =========================
// request id
// =========================

func newRequestID() string {

	b := make([]byte, 8)

	if _, err := rand.Read(b); err != nil {

		return strconv.FormatInt(
			time.Now().UnixNano(),
			16,
		)
	}

	return hex.EncodeToString(b)
}

// =========================
// real client ip
// =========================

func getClientIP(
	r *http.Request,
) string {

	xff := strings.TrimSpace(
		r.Header.Get("X-Forwarded-For"),
	)

	if xff != "" {

		parts := strings.Split(
			xff,
			",",
		)

		ip := strings.TrimSpace(
			parts[0],
		)

		if ip != "" {
			return ip
		}
	}

	xri := strings.TrimSpace(
		r.Header.Get("X-Real-IP"),
	)

	if xri != "" {
		return xri
	}

	host, _, err := net.SplitHostPort(
		r.RemoteAddr,
	)

	if err == nil && host != "" {
		return host
	}

	return r.RemoteAddr
}

// =========================
// masking
// =========================

func maskEmail(
	email string,
) string {

	if email == "" {
		return ""
	}

	parts := strings.Split(email, "@")

	if len(parts) != 2 {
		return "***"
	}

	name := parts[0]

	if len(name) <= 2 {
		return "***@" + parts[1]
	}

	return name[:2] + "***@" + parts[1]
}

func maskSub(
	sub string,
) string {

	if sub == "" {
		return ""
	}

	if len(sub) <= 8 {
		return "***"
	}

	return sub[:8] + "***"
}

// =========================
// structured log
// =========================

type requestLog struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`

	RequestID string `json:"request_id"`

	Sub   string `json:"sub,omitempty"`
	Email string `json:"email,omitempty"`

	Method string `json:"method"`
	Path   string `json:"path"`

	Status        int   `json:"status"`
	DurationMS    int64 `json:"duration_ms"`
	ResponseBytes int64 `json:"response_bytes"`

	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
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

		// =========================
		// request id context
		// =========================

		ctx := SetRequestID(
			r.Context(),
			requestID,
		)

		r = r.WithContext(ctx)

		rw := &responseWriter{
			ResponseWriter: w,
		}

		rw.Header().Set(
			"X-Request-ID",
			requestID,
		)

		next.ServeHTTP(
			rw,
			r,
		)

		// =========================
		// fallback status
		// =========================

		if rw.statusCode == 0 {
			rw.statusCode = http.StatusOK
		}

		authUser, _ := auth.GetAuthUser(
			r.Context(),
		)

		entry := requestLog{
			Timestamp: time.Now().
				UTC().
				Format(time.RFC3339),

			Level: "INFO",

			RequestID: requestID,

			Sub: maskSub(
				authUser.Sub,
			),

			Email: maskEmail(
				authUser.Email,
			),

			Method: r.Method,
			Path:   r.URL.Path,

			Status: rw.statusCode,

			DurationMS: time.Since(
				start,
			).Milliseconds(),

			ResponseBytes: rw.bytesWritten,

			ClientIP: getClientIP(
				r,
			),

			UserAgent: r.UserAgent(),
		}

		b, err := json.Marshal(
			entry,
		)

		if err != nil {

			log.Printf(
				`{"level":"ERROR","message":"failed to marshal request log"}`,
			)

			return
		}

		log.Println(
			string(b),
		)
	})
}
