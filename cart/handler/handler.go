package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"cart-go/helper/response"
	"cart-go/model"
	"cart-go/service"
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

func (h *handler) Get(ctx *gin.Context) {
	res, err := h.svc.Get()
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, "", err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "", res)
}

func (h *handler) GetDetail(ctx *gin.Context) {
	idStr, ok := ctx.GetQuery("cart_id")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, "", fmt.Errorf("query param cart_id should not be empty"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, "", err)
		return
	}

	res, err := h.svc.GetDetail(id)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, "", err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "", res)
}

func (h *handler) Create(ctx *gin.Context) {
	req := []model.CartRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, "", err)
		return
	}

	res, err := h.svc.Create(req)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, "", err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "", res)
}

func (h *handler) Delete(ctx *gin.Context) {
	idStr, ok := ctx.GetQuery("cart_id")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, "", fmt.Errorf("query param cart_id should not be empty"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, "", err)
		return
	}

	if err = h.svc.Delete(id); err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, "", err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "", nil)
}