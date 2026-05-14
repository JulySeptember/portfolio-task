package repository

import (
	"context"
	"time"
)

const queryTimeout = 5 * time.Second

func withTimeout(
	parent context.Context,
) (context.Context, context.CancelFunc) {

	return context.WithTimeout(
		parent,
		queryTimeout,
	)
}
