// internal/middleware/request_id.go

package middleware

import "context"

// =========================
// context key
// =========================

type contextKey string

const requestIDContextKey contextKey = "request_id"

// =========================
// set request id
// =========================

func SetRequestID(
	ctx context.Context,
	requestID string,
) context.Context {

	return context.WithValue(
		ctx,
		requestIDContextKey,
		requestID,
	)
}

// =========================
// get request id
// =========================

func GetRequestID(
	ctx context.Context,
) string {

	v, ok := ctx.Value(
		requestIDContextKey,
	).(string)

	if !ok {
		return ""
	}

	return v
}
