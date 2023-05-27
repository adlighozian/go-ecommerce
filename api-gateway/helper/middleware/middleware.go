package middleware

import (
	"api-gateway-go/helper/authjwt"
	"api-gateway-go/helper/response"
	"api-gateway-go/service"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
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
			Str("client_id", ctx.ClientIP()).
			Str("request_id", requestid.Get(ctx)).
			Str("user_agent", ctx.Request.UserAgent()).
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

// note: request counter can be rate limiter ;D
func RequestCounter(redis *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxBG := context.Background()
		clientIP := ctx.ClientIP()

		errIncr := redis.Incr(ctxBG, clientIP).Err()
		if errIncr != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Get the updated request count
		count, err := redis.Get(ctxBG, clientIP).Int64()
		if err != nil {
			// Handle the error (e.g., log, return an error response)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Header("X-Request-Count", strconv.FormatInt(count, 10))

		ctx.Next()
	}
}

func AuthMiddleware(allowedPaths []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uri string

		urlCtx, _ := ctx.Get("url")
		uri, _ = urlCtx.(string)

		parsedURL, _ := url.Parse(uri)
		path := parsedURL.Path[1:]

		// If path is the allowed, next to handler
		for _, allowedPath := range allowedPaths {
			if path == allowedPath {
				// Pass the real url and raw query into to handler
				ctx.Set("url", uri)
				ctx.Next()
				return
			}
		}

		// If the path is not allowed, authenticate the request
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.NewJSONResErr(ctx, http.StatusUnauthorized, "", "Missing authorization header")
			ctx.Abort()
			return
		}

		secretKey := viper.GetString("JWT_SECRET_KEY")

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.ParseWithClaims(
			tokenString,
			&authjwt.CustomClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			},
			jwt.WithIssuer("ecommerce"),
			jwt.WithIssuedAt(),
			jwt.WithLeeway(5*time.Second),
		)

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				response.NewJSONResErr(ctx, http.StatusUnauthorized, "", "Invalid token signature")
				ctx.Abort()
				return
			}
			response.NewJSONResErr(ctx, http.StatusBadRequest, "", "Invalid token")
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(*authjwt.CustomClaims)
		if !ok && !token.Valid {
			response.NewJSONResErr(ctx, http.StatusUnauthorized, "", "Invalid token claims")
			ctx.Abort()
			return
		}

		// Pass the userID, userRole to handler
		ctx.Set("userID", claims.UserID)
		ctx.Set("userRole", claims.UserRole)
		// Pass the real url and raw query into to handler
		ctx.Set("url", uri)
		ctx.Next()
	}
}

func HashedURLConverter(svc service.ShortenServiceI) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hashedURL := ctx.Param("hash")
		if hashedURL == "" {
			response.NewJSONResErr(ctx, http.StatusBadRequest, "", "")
			ctx.Abort()
			return
		}

		apiManagement, errSvc := svc.Get(hashedURL)
		if errSvc != nil {
			if errors.Is(errSvc, sql.ErrNoRows) {
				response.NewJSONResErr(ctx, http.StatusNotFound, "", "")
				ctx.Abort()
				return
			}
			_ = ctx.Error(errSvc)
			response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errSvc.Error())
			ctx.Abort()
			return
		}

		url := apiManagement.EndpointURL
		rawQuery := ctx.Request.URL.RawQuery
		if rawQuery != "" {
			url = url + "?" + rawQuery
		}

		// Pass the real url and raw query into to the next middleware
		ctx.Set("url", url)
		ctx.Next()
	}
}
