package middleware

import (
	"errors"
	"log"
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

				// =========================
				// broken pipe
				// =========================

				msg := strings.ToLower(
					err.Error(),
				)

				if strings.Contains(
					msg,
					"broken pipe",
				) ||
					strings.Contains(
						msg,
						"connection reset by peer",
					) {

					return
				}
			}

			// =========================
			// request id
			// =========================

			requestID := GetRequestID(
				r.Context(),
			)

			// =========================
			// response header
			// =========================

			if requestID != "" {

				w.Header().Set(
					"X-Request-ID",
					requestID,
				)
			}

			// =========================
			// panic logging
			// =========================

			log.Printf(
				"[PANIC] request_id=%s panic=%v\n%s",
				requestID,
				rec,
				debug.Stack(),
			)

			// =========================
			// response
			// =========================

			httpx.WriteError(
				w,
				http.StatusInternalServerError,
				httpx.CodeInternalServerError,
				"internal server error",
			)
		}()

		next.ServeHTTP(
			w,
			r,
		)
	})
}
