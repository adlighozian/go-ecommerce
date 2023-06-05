package handler

import (
	"size-go/helper/failerror"
	"size-go/helper/response"
	"size-go/model"
	"size-go/service"
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

func (h *handler) GetSize(ctx *gin.Context) {
	res, err := h.svc.GetSize()
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}

}

func (h *handler) CreateSize(ctx *gin.Context) {
	var data []model.SizeReq

	err := ctx.ShouldBindJSON(&data)
	failerror.FailError(err, "error bind json")

	res, err := h.svc.CreateSize(data)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}

func (h *handler) DeleteSize(ctx *gin.Context) {
	id := ctx.Query("size_id")
	var numi int

	if id != "" {
		num, err := strconv.Atoi(id)
		failerror.FailError(err, "error convert to int")
		numi = num
	}

	res, err := h.svc.DeleteSize(numi)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}
