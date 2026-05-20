// internal/handlers/timeout.go

package handlers

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultHandlerTimeout = 5 * time.Second
	listHandlerTimeout    = 10 * time.Second
)

// =========================
// withTimeout
// =========================

func withTimeout(
	r *http.Request,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {

	if timeout <= 0 {
		timeout = defaultHandlerTimeout
	}

	return context.WithTimeout(
		r.Context(),
		timeout,
	)
}
