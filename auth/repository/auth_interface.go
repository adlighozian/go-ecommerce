package repository

import (
	"auth-go/model"
	"time"
)

type AuthRepositoryI interface {
	Create(user *model.User) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	FirstOrCreate(user *model.User) (*model.User, error)

	SetRefreshToken(refreshToken string, dataByte []byte, refreshTokenDur time.Duration) error
	GetByRefreshToken(token string) (*model.RefreshToken, error)
}
