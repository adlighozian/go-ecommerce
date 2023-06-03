package handler

import (
	"order-go/helper/failerror"
	"order-go/helper/response"
	"order-go/model"
	"order-go/service"
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

func (h *handler) GetOrders(ctx *gin.Context) {

	idUser := ctx.Query("user_id")
	var numi int

	if idUser != "" {
		num, err := strconv.Atoi(idUser)
		failerror.FailError(err, "error convert to int")
		numi = num
	}

	res, err := h.svc.GetOrders(numi)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}

}

func (h *handler) GetOrdersByStoreID(ctx *gin.Context) {

	storeID := ctx.Query("store_id")
	var numi int

	if storeID != "" {
		num, err := strconv.Atoi(storeID)
		failerror.FailError(err, "error convert to int")
		numi = num
	}

	res, err := h.svc.GetOrdersByStoreID(numi)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}

}

func (h *handler) ShowOrders(ctx *gin.Context) {
	idUser := ctx.Query("user_id")
	orderNumber := ctx.Query("order_number")

	var numi int
	if idUser != "" {
		num, err := strconv.Atoi(idUser)
		failerror.FailError(err, "error convert to int")
		numi = num
	}

	var data model.OrderItems = model.OrderItems{
		UserId:      numi,
		OrderNumber: orderNumber,
	}

	res, err := h.svc.ShowOrders(data)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}

func (h *handler) CreateOrders(ctx *gin.Context) {
	var data model.GetOrders

	err := ctx.ShouldBindJSON(&data)
	failerror.FailError(err, "error bind json")

	res, err := h.svc.CreateOrders(data)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}

}

func (h *handler) UpdateOrders(ctx *gin.Context) {
	var data model.OrderUpd

	err := ctx.ShouldBindJSON(&data)
	failerror.FailError(err, "error bind json")

	res, err := h.svc.UpdateOrders(data)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}
