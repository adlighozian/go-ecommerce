package handler

import (
	"product-colors-go/helper/failerror"
	"product-colors-go/helper/response"
	"product-colors-go/model"
	"product-colors-go/service"
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

func (h *handler) GetColors(ctx *gin.Context) {
	res, err := h.svc.GetColors()
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}

}

func (h *handler) CreateColors(ctx *gin.Context) {
	var data []model.ColorsReq

	err := ctx.ShouldBindJSON(&data)
	failerror.FailError(err, "error bind json")

	res, err := h.svc.CreateColors(data)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}

func (h *handler) DeleteColors(ctx *gin.Context) {
	id := ctx.Query("color_id")
	var numi int

	if id != "" {
		num, err := strconv.Atoi(id)
		failerror.FailError(err, "error convert to int")
		numi = num
	}

	res, err := h.svc.DeleteColors(numi)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}
