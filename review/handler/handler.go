package handler

import (
	"fmt"
	"net/http"
	"review-go/helper/response"
	"review-go/model"
	"review-go/service"
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

func (h *handler) GetByProductID(ctx *gin.Context) {
	productIDString, ok := ctx.GetQuery("product_id")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param Review_id should not be empty"))
		return
	}

	productID, err := strconv.Atoi(productIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	if productID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("produtID should be positive number"))
		return
	}

	res, err := h.svc.GetByProductID(productID)
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

	for _, v := range req {
		if v.UserID <= 0 {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("user_id should be positive number"))
			return
		}

		if v.ProductID <= 0 {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("product_id should be positive number"))
			return
		}

		if v.Rating <= 0 {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("rating should not be empty"))
			return
		}

		if v.ReviewText == "" {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("review_text should not be empty"))
			return
		}
	}

	res, err := h.svc.Create(req)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}