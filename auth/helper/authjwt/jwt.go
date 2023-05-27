package authjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates an access token for a given user ID
func GenerateAccessToken(userID uint, userRole string) (string, error) {
	secretKey := viper.GetString("JWT_SECRET_KEY")
	accessTokenDuration := viper.GetDuration("JWT_ACCESS_TOKEN_DURATION")

	claims := &CustomClaims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ecommerce",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// GenerateRefreshToken generates a refresh token for a given user ID
func GenerateRefreshToken(userID uint, userRole string) (string, error) {
	secretKey := viper.GetString("JWT_SECRET_KEY")
	refreshTokenDuration := viper.GetDuration("JWT_REFRESH_TOKEN_DURATION")

	claims := &CustomClaims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ecommerce",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
