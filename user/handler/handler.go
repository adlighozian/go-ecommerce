package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"auth-go/service"
)

type handler struct {
	svc service.Servicer
}

func NewHandler(svc service.Servicer) *Handlerer {
	return &handler{
		svc: svc,
	}
}

func (h *handler) Get(ctx *gin.Context) {
	h.svc.GetList()
}

func (h *handler) GetDetail(ctx *gin.Context) {
	h.svc.GetDetail()
}

func (h *handler) Create(ctx *gin.Context) {
	h.svc.Create()
}

func (h *handler) Update(ctx *gin.Context) {
	h.svc.Update()
}

func (h *handler) Delete(ctx *gin.Context) {
	h.svc.Delete()
}