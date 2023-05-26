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
	userIDString, ok := ctx.GetQuery("user_ID")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param user_id should not be empty"))
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := h.svc.Get(userID)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) GetDetail(ctx *gin.Context) {
	userIDString, ok := ctx.GetQuery("user_ID")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param user_id should not be empty"))
		return
	}

	cartIDString, ok := ctx.GetQuery("cart_id")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param cart_id should not be empty"))
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	cartID, err := strconv.Atoi(cartIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := h.svc.GetDetail(userID, cartID)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) Create(ctx *gin.Context) {
	req := []model.CartRequest{}

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

func (h *handler) Delete(ctx *gin.Context) {
	userIDString, ok := ctx.GetQuery("user_ID")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param user_id should not be empty"))
		return
	}

	cartIDString, ok := ctx.GetQuery("cart_id")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param cart_id should not be empty"))
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	cartID, err := strconv.Atoi(cartIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.svc.Delete(userID, cartID); err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, nil)
}