package timeout

import (
	"context"
	"time"
)

func NewCtxTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
