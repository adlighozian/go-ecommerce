package handler

import (
	"product-go/service"

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

func (h *handler) GetProduct(ctx *gin.Context) {


	
}

func (h *handler) ShowProduct(ctx *gin.Context) {

}

func (h *handler) CreateProduct(ctx *gin.Context) {

}

func (h *handler) UpdateProduct(ctx *gin.Context) {

}

func (h *handler) DeleteProduct(ctx *gin.Context) {

}
