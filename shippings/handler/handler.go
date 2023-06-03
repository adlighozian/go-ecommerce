package handler

import (
	"shippings-go/helper/failerror"
	"shippings-go/helper/response"
	"shippings-go/model"
	"shippings-go/service"
	"strconv"

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

func (h *handler) GetShipping(ctx *gin.Context) {
	res, err := h.svc.GetShipping()
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}

}

func (h *handler) CreateShipping(ctx *gin.Context) {
	var data []model.ShippingReq

	err := ctx.ShouldBindJSON(&data)
	failerror.FailError(err, "error bind json")

	res, err := h.svc.CreateShipping(data)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}

func (h *handler) DeleteShipping(ctx *gin.Context) {
	id := ctx.Query("shipping_id")
	var numi int

	if id != "" {
		num, err := strconv.Atoi(id)
		failerror.FailError(err, "error convert to int")
		numi = num
	}

	res, err := h.svc.DeleteShipping(numi)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}
