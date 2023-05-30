package service

import (
	"auth-go/helper/authjwt"
	"auth-go/model"
	"auth-go/repository"
	"errors"
	"time"
)

type AuthService struct {
	repo repository.AuthRepositoryI
}

func NewAuthService(repo repository.AuthRepositoryI) AuthServiceI {
	svc := new(AuthService)
	svc.repo = repo
	return svc
}

func (svc *AuthService) Create(registerReq *model.RegisterReq) (*model.User, error) {
	hashedPassword, errHash := authjwt.HashPassword(registerReq.Password)
	if errHash != nil {
		return nil, errHash
	}

	newUser := &model.User{
		Username: registerReq.Username,
		Email:    registerReq.Email,
		Password: hashedPassword,
		Role:     "user",
	}
	return svc.repo.Create(newUser)
}

func (svc *AuthService) FirstOrCreate(userReq *model.UserReq) (*model.User, error) {
	newUser := &model.User{
		Username: userReq.Username,
		Email:    userReq.Email,
		Role:     "user",
		FullName: userReq.FullName,
		ImageURL: userReq.ImageURL,
	}
	return svc.repo.FirstOrCreate(newUser)
}

func (svc *AuthService) GetByEmail(loginReq *model.LoginReq) (*model.User, error) {
	user, errRepo := svc.repo.GetByEmail(loginReq.Email)
	if errRepo != nil {
		return nil, errRepo
	}

	if isPasswordCorrect := authjwt.CheckPasswordHash(loginReq.Password, user.Password); !isPasswordCorrect {
		return nil, errors.New("incorrect password")
	}

	return user, nil
}

func (svc *AuthService) SetRefreshToken(refreshToken string, dataByte []byte, refreshTokenDur time.Duration) error {
	return svc.repo.SetRefreshToken(refreshToken, dataByte, refreshTokenDur)
}

func (svc *AuthService) GetByRefreshToken(token string) (*model.RefreshToken, error) {
	return svc.repo.GetByRefreshToken(token)
}
