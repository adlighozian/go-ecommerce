package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Logger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		ctx.Next()

		latency := time.Since(start)
		statusCode := ctx.Writer.Status()

		msg := "Request"
		if len(ctx.Errors) > 0 {
			msg = ctx.Errors.String()
		}

		logger := logger.With().
			Str("method", ctx.Request.Method).
			Str("path", path).
			// Str("client_id", ctx.ClientIP()).
			// Str("user_agent", ctx.Request.UserAgent()).
			Int("status_code", statusCode).
			Int("body_size", ctx.Writer.Size()).
			Dur("latency", latency).
			Logger()

		switch {
		case statusCode >= http.StatusInternalServerError:
			logger.WithLevel(zerolog.ErrorLevel).Msg(msg)
		case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
			logger.WithLevel(zerolog.WarnLevel).Msg(msg)
		default:
			logger.WithLevel(zerolog.InfoLevel).Msg(msg)
		}
	}
}
