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

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
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

		var uri string
		urlCtx, exists := ctx.Get("url")
		if exists {
			uri, _ = urlCtx.(string)
		}

		logger := logger.With().
			Str("method", ctx.Request.Method).
			Str("hashed_path", path).
			Str("client_id", ctx.ClientIP()).
			Str("request_id", requestid.Get(ctx)).
			Str("endpoint_url", uri).
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

		// Increment the request count
		errIncr := redis.Incr(ctxBG, clientIP).Err()
		if errIncr != nil {
			_ = ctx.Error(errIncr)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Get the updated request count
		count, errGet := redis.Get(ctxBG, clientIP).Int64()
		if errGet != nil {
			_ = ctx.Error(errGet)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Header("X-Request-Count", strconv.FormatInt(count, 10))

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

		// Pass the real url and raw query into to the AuthMiddleware
		ctx.Set("url", url)
		ctx.Next()
	}
}

// If path is the allowed, next to handler.
func isAllowedPath(ctx *gin.Context, allowedPaths []string) (string, string, bool) {
	var uri string

	urlCtx, _ := ctx.Get("url")
	uri, _ = urlCtx.(string)

	parsedURL, _ := url.Parse(uri)
	path := parsedURL.Path[1:]

	for _, allowedPath := range allowedPaths {
		if path == allowedPath {
			// Pass the real url and raw query into to handler
			ctx.Set("url", uri)
			ctx.Next()
			return "", "", true
		}
	}
	return uri, path, false
}

func AuthMiddleware(jwtSecretKey string, allowedPaths []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri, path, shouldReturn := isAllowedPath(ctx, allowedPaths)
		if shouldReturn {
			return
		}

		// If the path is not allowed, authenticate the request
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.NewJSONResErr(ctx, http.StatusUnauthorized, "", "Missing authorization header")
			ctx.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, errParse := jwt.ParseWithClaims(
			tokenString,
			&authjwt.CustomClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecretKey), nil
			},
			jwt.WithIssuer("ecommerce"),
			jwt.WithIssuedAt(),
			jwt.WithLeeway(5*time.Second),
		)

		if errParse != nil {
			if errors.Is(errParse, jwt.ErrSignatureInvalid) {
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

		// Pass the userID, userRole to AuthzMiddleware
		ctx.Set("userID", claims.UserID)
		ctx.Set("userRole", claims.UserRole)
		// Pass the real url, raw query, and path into to AuthzMiddleware
		ctx.Set("url", uri)
		ctx.Set("path", path)
		ctx.Next()
	}
}

func AuthzMiddleware(enforcer *casbin.Enforcer, allowedPaths []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri, path, shouldReturn := isAllowedPath(ctx, allowedPaths)
		if shouldReturn {
			return
		}

		var userID string
		userIDCtx, _ := ctx.Get("userID")
		userID, _ = userIDCtx.(string)

		var userRole string
		userRoleCtx, _ := ctx.Get("userRole")
		userRole, _ = userRoleCtx.(string)

		// Now check if the user has the necessary permissions
		sub := userRole
		obj := path               // The object the user is trying to access
		act := ctx.Request.Method // The action the user is trying to perform

		errLoad := enforcer.LoadPolicy()
		if errLoad != nil {
			// Handle error
			_ = ctx.Error(errLoad)
			response.NewJSONResErr(ctx, http.StatusInternalServerError, "", "failed to load policy")
			ctx.Abort()
			return
		}

		ok, errEnforce := enforcer.Enforce(sub, obj, act)
		if errEnforce != nil {
			// Handle error
			_ = ctx.Error(errEnforce)
			response.NewJSONResErr(ctx, http.StatusInternalServerError, "", "unable to enforce policy")
			ctx.Abort()
			return
		}

		if !ok {
			// The user is not authorized, return an error
			response.NewJSONResErr(ctx, http.StatusForbidden, "", "You are not authorized to access this resource")
			ctx.Abort()
			return
		}

		// Pass the userID to handler
		ctx.Set("userID", userID)
		// Pass the real url and raw query into to handler
		ctx.Set("url", uri)
		ctx.Next()
	}
}
