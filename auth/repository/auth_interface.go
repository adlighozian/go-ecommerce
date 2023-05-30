package repository

import (
	"auth-go/model"
	"time"
)

type AuthRepositoryI interface {
	// Get(..., err error)
	// GetDetail(..., err error)
	Create(user *model.User) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	FirstOrCreate(user *model.User) (*model.User, error)
	// Update(..., err error)
	// Delete(..., err error)

	SetRefreshToken(refreshToken string, dataByte []byte, refreshTokenDur time.Duration) error
	GetByRefreshToken(token string) (*model.RefreshToken, error)
}
