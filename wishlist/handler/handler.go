package handler

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"wishlist-go/helper/response"
	"wishlist-go/model"
	"wishlist-go/service"
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
	idStr, ok := ctx.GetQuery("wishlist_id")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, "", fmt.Errorf("query param wishlist_id should not be empty"))
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
	req := []model.WishlistRequest{}

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
	idStr, ok := ctx.GetQuery("wishlist_id")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, "", fmt.Errorf("query param wishlist_id should not be empty"))
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