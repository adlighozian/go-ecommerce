package timeout

import (
	"context"
	"time"
)

func NewCtxTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 50*time.Second)
}