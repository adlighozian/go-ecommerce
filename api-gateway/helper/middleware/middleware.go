package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Logger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		msg := "Request"
		if len(c.Errors) > 0 {
			msg = c.Errors.String()
		}

		logger := logger.With().
			Str("method", c.Request.Method).
			Str("path", path).
			Str("client_id", c.ClientIP()).
			Str("request_id", requestid.Get(c)).
			Str("user_agent", c.Request.UserAgent()).
			Int("status_code", statusCode).
			Int("body_size", c.Writer.Size()).
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
