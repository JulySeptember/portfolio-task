// internal/middleware/recovery.go

package middleware

import (
	"errors"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"

	"portfolio/backend/internal/httpx"
)

func Recovery(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		defer func() {

			rec := recover()

			if rec == nil {
				return
			}

			// =========================
			// broken pipe / connection reset
			// =========================

			var err error

			switch v := rec.(type) {

			case error:
				err = v

			case string:
				err = errors.New(v)
			}

			if err != nil {

				// net.Error
				var netErr net.Error

				if errors.As(err, &netErr) &&
					!netErr.Timeout() {

					return
				}

				// broken pipe
				msg := strings.ToLower(err.Error())

				if strings.Contains(msg, "broken pipe") ||
					strings.Contains(msg, "connection reset by peer") {

					return
				}
			}

			// =========================
			// panic logging
			// =========================

			log.Printf(
				"[PANIC] %v\n%s",
				rec,
				debug.Stack(),
			)

			// =========================
			// response
			// =========================

			httpx.WriteError(
				w,
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				"internal server error",
			)
		}()

		next.ServeHTTP(w, r)
	})
}
