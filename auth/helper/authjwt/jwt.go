package authjwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a token for a given token duration, user ID, and user Role.
func GenerateToken(jwtSecretKey string, tokenDur time.Duration, userID uint, userRole string) (string, error) {
	claims := &CustomClaims{
		UserID:   strconv.FormatUint(uint64(userID), 10),
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ecommerce",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDur)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}
