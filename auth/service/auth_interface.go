package service

import (
	"auth-go/model"
	"time"
)

type AuthServiceI interface {
	Create(registerReq *model.RegisterReq) (*model.User, error)
	FirstOrCreate(provider string, userReq *model.UserReq) (*model.User, error)
	LoginByEmail(loginReq *model.LoginReq) (*model.User, error)

	SetRefreshToken(refreshToken string, dataByte []byte, refreshTokenDur time.Duration) error
	GetByRefreshToken(token string) (*model.RefreshToken, error)
}
