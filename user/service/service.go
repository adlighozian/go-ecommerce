package  service

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"auth-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) *Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) Get(ctx *gin.Context) {
	svc.repo.Get()
}

func (svc *service) GetDetail(ctx *gin.Context) {
	svc.repo.GetDetail()
}

func (svc *service) Create(ctx *gin.Context) {
	svc.repo.Create()
}

func (svc *service) Update(ctx *gin.Context) {
	svc.repo.Update()
}

func (svc *service) Delete(ctx *gin.Context) {
	svc.repo.Delete()
}