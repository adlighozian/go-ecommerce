package handler

import (
	"product-go/helper/failerror"
	"product-go/helper/response"
	"product-go/model"
	"product-go/service"
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

func (h *handler) GetProduct(ctx *gin.Context) {
	brand := ctx.Query("brand")
	category := ctx.Query("category")
	name := ctx.Query("name")

	data := model.ProductSearch{
		Brand:    brand,
		Category: category,
		Name:     name,
	}

	res, err := h.svc.GetProduct(data)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}

}

func (h *handler) ShowProduct(ctx *gin.Context) {
	id := ctx.Query("id")
	var numi int

	if id != "" {
		num, err := strconv.Atoi(id)
		failerror.FailError(err, "error convert to int")
		numi = num
	}

	res, err := h.svc.ShowProduct(numi)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}

func (h *handler) CreateProduct(ctx *gin.Context) {
	var data []model.ProductReq

	err := ctx.ShouldBindJSON(&data)
	failerror.FailError(err, "error bind json")

	res, err := h.svc.CreateProduct(data)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}

func (h *handler) UpdateProduct(ctx *gin.Context) {
	var data model.ProductReq
	err := ctx.ShouldBindJSON(&data)
	failerror.FailError(err, "error bind json")

	res, err := h.svc.UpdateProduct(data)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}

func (h *handler) DeleteProduct(ctx *gin.Context) {
	id := ctx.Query("id")
	var numi int

	if id != "" {
		num, err := strconv.Atoi(id)
		failerror.FailError(err, "error convert to int")
		numi = num
	}

	res, err := h.svc.DeleteProduct(numi)
	if err != nil {
		response.ResponseError(ctx, res.Status, err)
	} else {
		response.ResponseSuccess(ctx, res.Status, res.Data)
	}
}
