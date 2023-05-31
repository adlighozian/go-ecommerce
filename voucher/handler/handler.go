package handler

import (
	"strconv"
	"voucher-go/helper/failerror"
	"voucher-go/helper/response"
	"voucher-go/model"
	"voucher-go/service"

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

func (h *handler) GetVoucher(ctx *gin.Context) {
	id := ctx.Query("user_id")
	var numi int

	if id != "" {
		num, err := strconv.Atoi(id)
		failerror.FailError(err, "error convert to int")
		numi = num
	}

	res, err := h.svc.GetVoucher(numi)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}

}

func (h *handler) ShowVoucher(ctx *gin.Context) {
	code := ctx.Query("code")

	res, err := h.svc.ShowVoucher(code)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}

func (h *handler) CreateVoucher(ctx *gin.Context) {
	var data []model.VoucherReq

	err := ctx.ShouldBindJSON(&data)
	failerror.FailError(err, "error bind json")

	res, err := h.svc.CreateVoucher(data)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}

func (h *handler) DeleteVoucher(ctx *gin.Context) {
	id := ctx.Query("id")
	var numi int

	if id != "" {
		num, err := strconv.Atoi(id)
		failerror.FailError(err, "error convert to int")
		numi = num
	}

	res, err := h.svc.DeleteVoucher(numi)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}
