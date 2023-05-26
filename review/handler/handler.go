package handler

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"review-go/helper/response"
	"review-go/model"
	"review-go/service"
)

type handler struct {
	svc service.Servicer
}

func NewHandler(svc service.Servicer) Handlerer {
	return &handler{
		svc: svc,
	}
}

func (h *handler) GetByProductID(ctx *gin.Context) {
	idStr, ok := ctx.GetQuery("Review_id")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param Review_id should not be empty"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := h.svc.GetByProductID(id)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) Create(ctx *gin.Context) {
	req := []model.ReviewRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := h.svc.Create(req)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}