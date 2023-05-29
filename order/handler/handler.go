package handler

import (
	"order-go/service"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc service.Servicer
}

func NewHandler(svc service.Servicer) Handlerer {
	return &handler{
		svc: svc,
	}
}

func (h *handler) GetOrders(ctx *gin.Context) {

}

func (h *handler) ShowOrders(ctx *gin.Context) {

}

func (h *handler) CreateOrders(ctx *gin.Context) {

}

func (h *handler) UpdateOrders(ctx *gin.Context) {

}
